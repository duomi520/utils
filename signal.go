package utils

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Parent interface {
	Get() any
	String() string
	Int() int
}

// Signal 信号-状态单向传播
// 一个数据的容器，当它存储的数据改变时，依赖于这个 Signal 的计算函数Computer标注“脏” 状态、副作用Effector可以自动运行
// 采用所谓的 “推后拉” 模型：“推” 阶段，在 Signal 变为 “脏”（即其值发生了改变）时，会递归地把 “脏” 状态传递到依赖它的所有计算函数Computer上，所有潜在的重新计算都被推迟，直到显式地 Operate 某个 Computer 的值。

type Signal struct {
	index int
	u     *Universe
}

// Get 线程不安全,需控制在evaluate函数内运行
func (s Signal) Get() any {
	return s.u.signalValue[s.index]
}

// String 线程不安全,需控制在evaluate函数内运行
func (s Signal) String() string {
	return s.u.signalValue[s.index].(string)
}

// Int 线程不安全,需控制在evaluate函数内运行
func (s Signal) Int() int {
	return s.u.signalValue[s.index].(int)
}

// Computer 衍生 - 衍生能缓存计算结果，避免重复的计算
// 惰性求值（lazy evaluate）- 只有被使用到的才会计算结果
type Computer struct {
	u *Universe
	//所有的Signal父代
	parentSignal []int
	//求值链
	evaluateChain []*Computer
	//求值函数
	evaluate func() any
	//Reactions 反应 - 反应是数据更新时的监听器，监视值修改后，立即执行
	effect func(any)
	//false-需计算值，true-无需计算值
	renovate bool
	value    any
}

// NewComputer 线程不安全
func NewComputer(u *Universe, effect func(any), warp func() (func() any, []Parent)) *Computer {
	evaluate, p := warp()
	c := &Computer{
		u:        u,
		evaluate: evaluate,
		effect:   effect,
		renovate: false,
	}
	for _, v := range p {
		switch p := v.(type) {
		case Signal:
			u.signalChild[p.index] = append(u.signalChild[p.index], c)
			c.parentSignal = append(c.parentSignal, p.index)
			if p.u != u {
				panic("NewComputer: Signal不同Universe")
			}
		case *Computer:
			if p.u != u {
				panic("NewComputer: Computer不同Universe")
			}
			c.parentSignal = append(c.parentSignal, p.parentSignal...)
			for _, i := range c.parentSignal {
				u.signalChild[i] = append(u.signalChild[i], c)
			}
			c.evaluateChain = append(c.evaluateChain, p.evaluateChain...)
			c.evaluateChain = append(c.evaluateChain, p)
		default:
			panic(fmt.Sprintf("NewComputer: 无效的参数类型:%v", v))
		}
	}
	c.parentSignal = removeDuplicates(c.parentSignal)
	c.evaluateChain = uniqueWithoutSorting(c.evaluateChain)
	return c
}

// removeDuplicates 排序去重
func removeDuplicates(s []int) []int {
	if len(s) == 0 {
		return nil
	}
	// 先排序
	sort.Ints(s)
	// k用来记录不重复元素的索引位置
	k := 0
	for i := 1; i < len(s); i++ {
		if s[k] != s[i] {
			k++
			// 将不重复的元素移到数组前面
			s[k] = s[i]
		}
	}
	// 返回不重复元素的部分切片
	return s[:k+1]
}

// uniqueWithoutSorting 去重但不排序
func uniqueWithoutSorting(s []*Computer) []*Computer {
	if len(s) == 0 {
		return nil
	}
	// 使用map来记录元素是否已经出现过
	seen := make(map[*Computer]bool)
	var result []*Computer
	for _, value := range s {
		// 如果这个值没有在map中出现过，就添加到结果切片中
		if _, ok := seen[value]; !ok {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}

// Get 线程不安全,需控制在evaluate函数内运行
func (c *Computer) Get() any {
	if c.renovate {
		return c.value
	}
	for _, v := range c.evaluateChain {
		if !v.renovate {
			v.eval()
		}
	}
	c.eval()
	return c.value
}

func (c *Computer) eval() {
	c.value = c.evaluate()
	c.renovate = true
	if c.effect != nil {
		c.effect(c.value)
	}
}

// String 线程不安全,需控制在evaluate函数内运行
func (c *Computer) String() string {
	return c.Get().(string)
}

// Int 线程不安全,需控制在evaluate函数内运行
func (c *Computer) Int() int {
	return c.Get().(int)
}

type setSignalMsg struct {
	index int
	val   any
}

type operateComputerMsg struct {
	c  *Computer
	do func(any)
}

// Universe 响应式状态管理,放置在一同协程中运算，规避锁的影响
// 适用于生产者慢，消费者快的场景

type Universe struct {
	signalValue         []any
	signalEffect        []func(any)
	signalChild         [][]*Computer
	setSignalChan       chan setSignalMsg
	operateComputerChan chan operateComputerMsg
	closeOnce           sync.Once
	// 0: stop 1: live
	stopFlag int32
	stopChan chan struct{}
}

func NewUniverse() *Universe {
	u := &Universe{
		setSignalChan:       make(chan setSignalMsg, 128),
		operateComputerChan: make(chan operateComputerMsg, 128),
		stopFlag:            1,
		stopChan:            make(chan struct{}),
	}
	return u
}

// Close 关闭
func (u *Universe) Close() {
	u.closeOnce.Do(func() {
		atomic.StoreInt32(&u.stopFlag, 0)
		close(u.stopChan)
		time.Sleep(3 * time.Second)
		u.signalValue = nil
		u.signalEffect = nil
		u.signalChild = nil
		close(u.setSignalChan)
		close(u.operateComputerChan)
	})
}

// 线程不安全
func (u *Universe) NewSignal(v any, effect func(any)) Signal {
	u.signalValue = append(u.signalValue, v)
	u.signalEffect = append(u.signalEffect, effect)
	u.signalChild = append(u.signalChild, make([]*Computer, 0, 16))
	return Signal{index: len(u.signalEffect) - 1, u: u}
}

// SetSignal 非原子性，乱序执行。
func (u *Universe) SetSignal(s Signal, a any) {
	if atomic.LoadInt32(&u.stopFlag) != 0 {
		u.setSignalChan <- setSignalMsg{s.index, a}
	}
}

// Operate 非原子性，乱序执行。
func (u *Universe) Operate(c *Computer, do func(any)) {
	if atomic.LoadInt32(&u.stopFlag) != 0 {
		u.operateComputerChan <- operateComputerMsg{c: c, do: do}
	}
}

func (u *Universe) Run() {
	if len(u.signalChild) > 0 {
		for i := range u.signalChild {
			u.signalChild[i] = uniqueWithoutSorting(u.signalChild[i])
		}
	}
	go func() {
		for {
			select {
			case msg := <-u.setSignalChan:
				u.signalValue[msg.index] = msg.val
				for _, v := range u.signalChild[msg.index] {
					v.renovate = false
				}
				if u.signalEffect[msg.index] != nil {
					u.signalEffect[msg.index](msg.val)
				}
			case msg := <-u.operateComputerChan:
				if msg.do != nil {
					msg.do(msg.c.Get())
				} else {
					msg.c.Get()
				}
			case <-u.stopChan:
				return
			}
		}
	}()
}

// https://zhuanlan.zhihu.com/p/691797618
// https://developer.aliyun.com/article/1218778
