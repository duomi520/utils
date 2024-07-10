package utils

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

type Topic interface {
	UnSubscribe(int64)
	Subscribe(int64, func())
}

// signal 一个数据的容器，当它存储的数据改变时，依赖于这个 Signal 的计算函数或者副作用可以自动更新
// 采用所谓的 “推后拉” 模型：“推” 阶段，在 Signal 变为 “脏”（即其值发生了改变）时，会递归地把 “脏” 状态传递到依赖它的所有 Signals 上，所有潜在的重新计算都被推迟，直到显式地请求某个 Signal 的值。
// 适用于生产者慢，消费者快的场景
// 惰性求值（lazy evaluate）- 只有被使用到的才会更新

type Signal struct {
	id       int64
	value    any
	subIndex []int64
	subSet   []func()
}

func (s *Signal) Subscribe(i int64, f func()) {
	s.subIndex = append(s.subIndex, i)
	s.subSet = append(s.subSet, f)
}

func (s *Signal) UnSubscribe(i int64) {
	index := slices.Index(s.subIndex, i)
	if index > -1 {
		s.subIndex = slices.Delete(s.subIndex, index, index+1)
		s.subSet = slices.Delete(s.subSet, index, index+1)
	}
}

type Computer struct {
	id int64
	//false-需求值，true-无需更新
	renovate bool
	value    any
	parent   []int64
	subIndex []int64
	subSet   []func()
	evaluate func(*Processor) any
}

func (c *Computer) Subscribe(i int64, f func()) {
	c.subIndex = append(c.subIndex, i)
	c.subSet = append(c.subSet, f)
}

func (c *Computer) UnSubscribe(i int64) {
	index := slices.Index(c.subIndex, i)
	if index > -1 {
		c.subIndex = slices.Delete(c.subIndex, index, index+1)
		c.subSet = slices.Delete(c.subSet, index, index+1)
	}
}

// command 指令
type command struct {
	id     int64
	value  any
	do     func(int64, any) any
	report chan any
}

// Processor 处理器 响应式状态管理,放置在一同协程中运算，规避锁的影响
type Processor struct {
	increases   int64
	set         map[int64]any
	commandChan chan command
	closeOnce   sync.Once
	// 0: stop 1: live
	stopFlag int32
	stopChan chan struct{}
}

func NewProcessor() *Processor {
	p := &Processor{
		set:         make(map[int64]any),
		commandChan: make(chan command, 128),
		stopFlag:    1,
		stopChan:    make(chan struct{}),
	}
	go func() {
		for {
			select {
			case c := <-p.commandChan:
				v := c.do(c.id, c.value)
				if c.report != nil {
					c.report <- v
				}
			case <-p.stopChan:
				return
			}
		}
	}()
	return p
}

// Close 停止
func (p *Processor) Close() {
	p.closeOnce.Do(func() {
		atomic.StoreInt32(&p.stopFlag, 0)
		close(p.stopChan)
		time.Sleep(5 * time.Second)
		p.set = nil
		close(p.commandChan)
	})

}

//------------------------------------------------------------------------------

// GetSignal
// 该函数仅可用于协程内部计算
func GetSignal(p *Processor, id int64) any {
	v, ok := p.set[id]
	if !ok {
		return nil
	}
	return v.(*Signal).value
}

// SetSignal
// 该函数仅可用于协程内部计算
func SetSignal(p *Processor, id int64, newValue any) {
	v, ok := p.set[id]
	if ok {
		s := v.(*Signal)
		s.value = newValue
		for _, execute := range s.subSet {
			execute()
		}
	}
}

// GetComputer
// 该函数仅用于协程内部计算
func GetComputer(p *Processor, id int64) any {
	v, ok := p.set[id]
	if !ok {
		return nil
	}
	c := v.(*Computer)
	if !c.renovate {
		c.value = c.evaluate(p)
		c.renovate = true
	}
	return c.value
}

//------------------------------------------------------------------------------

func (p *Processor) newSignal(id int64, newValue any) any {
	p.set[id] = &Signal{id: id, value: newValue}
	return nil
}

// SignalState 信号 状态单向传播
func (p *Processor) SignalState(newValue any) (id int64) {
	id = atomic.AddInt64(&p.increases, 1)
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return 0
	}
	p.commandChan <- command{id: id, value: newValue, do: p.newSignal}
	return
}

func (p *Processor) setSignal(id int64, newValue any) any {
	SetSignal(p, id, newValue)
	return nil
}

// SetSignalState 值变化，通知订阅者，状态变 “脏”，会阻塞,
func (p *Processor) SetSignalState(id int64, newValue any) {
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return
	}
	p.commandChan <- command{id: id, value: newValue, do: p.setSignal}
}

func (p *Processor) getSignal(id int64, newValue any) any {
	v, ok := p.set[id]
	if !ok {
		return nil
	}
	s := v.(*Signal)
	return s.value
}

// GetSignalState
func (p *Processor) GetSignalState(id int64) any {
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return nil
	}
	c := make(chan any)
	p.commandChan <- command{id: id, do: p.getSignal, report: c}
	r := <-c
	close(c)
	return r
}

func (p *Processor) computed(id int64, value any) any {
	v := value.(*Computer)
	p.set[id] = v
	//通知状态修改函数
	f := func() {
		v.renovate = false
		for _, execute := range v.subSet {
			execute()
		}
	}
	//订阅
	for _, sub := range v.parent {
		c, ok := p.set[sub]
		if ok {
			c.(Topic).Subscribe(id, f)
		}
	}
	return nil
}

// Computed 衍生 衍生能缓存计算结果，避免重复的计算，并且也能自动追踪依赖以及同步更新
func (p *Processor) Computed(do func(*Processor) any, t ...int64) (id int64) {
	id = atomic.AddInt64(&p.increases, 1)
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return 0
	}
	p.commandChan <- command{id: id, value: &Computer{id: id, renovate: false, parent: t, evaluate: do}, do: p.computed}
	return id
}

func (p *Processor) getComputer(id int64, v any) any {
	return GetComputer(p, id)
}

// GetComputerState 求值，状态为“新”时，返回缓存，“脏”时，拉取订阅者，重新计算，会阻塞，有缓存，取得值的不一定是最新的。
func (p *Processor) GetComputerState(id int64) any {
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return nil
	}
	c := make(chan any)
	p.commandChan <- command{id: id, do: p.getComputer, report: c}
	r := <-c
	close(c)
	return r
}

// RemoveComputer
func (p *Processor) removeComputer(id int64, value any) any {
	v, ok := p.set[id]
	if !ok {
		return errors.New("computer no found")
	}
	if len(v.(*Computer).subIndex) > 0 {
		return errors.New("sub is exist")
	}
	for _, parent := range v.(*Computer).parent {
		c, ok := p.set[parent]
		if ok {
			c.(Topic).UnSubscribe(id)
		}
	}
	delete(p.set, id)
	v.(*Computer).parent = nil
	return nil
}

// RemoveComputer
func (p *Processor) RemoveComputer(id int64) error {
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return errors.New("processor closed")
	}
	c := make(chan any)
	p.commandChan <- command{id: id, do: p.removeComputer, report: c}
	r := <-c
	close(c)
	if r == nil {
		return nil
	}
	return r.(error)
}

type effectorMsg struct {
	do     func()
	parent []int64
}

func (p *Processor) effector(id int64, value any) any {
	e := value.(effectorMsg)
	p.set[id] = e.parent
	//订阅
	for _, sub := range e.parent {
		c, ok := p.set[sub]
		if ok {
			c.(Topic).Subscribe(id, e.do)
		}
	}
	return nil
}

// Effector Reactions 反应 反应是数据更新时的监听器，监视值修改后，立即执行
func (p *Processor) Effector(do func(*Processor), t ...int64) (id int64) {
	id = atomic.AddInt64(&p.increases, 1)
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return 0
	}
	e := effectorMsg{do: func() { do(p) }, parent: t}
	p.commandChan <- command{id: id, value: e, do: p.effector}
	return
}

func (p *Processor) removeEffector(id int64, value any) any {
	v, ok := p.set[id]
	if !ok {
		return nil
	}
	for _, parent := range v.([]int64) {
		c, ok := p.set[parent]
		if ok {
			c.(Topic).UnSubscribe(id)
		}
	}
	delete(p.set, id)
	return nil
}

// RemoveEffector
func (p *Processor) RemoveEffector(id int64) {
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return
	}
	p.commandChan <- command{id: id, do: p.removeEffector}
}

func (p *Processor) unSubscribeEffector(id int64, parent any) any {
	v, ok := p.set[parent.(int64)]
	if !ok {
		return nil
	}
	v.(Topic).UnSubscribe(id)
	return nil
}

// UnSubscribeEffector
func (p *Processor) UnSubscribeEffector(id, parent int64) {
	if atomic.LoadInt32(&p.stopFlag) == 0 {
		return
	}
	p.commandChan <- command{id: id, value: parent, do: p.unSubscribeEffector}
}

// https://zhuanlan.zhihu.com/p/691797618
// https://developer.aliyun.com/article/1218778
