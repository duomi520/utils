package utils

import (
	"sync"
	"sync/atomic"
	"time"
)

// Signal 信号-状态单向传播
// 一个数据的容器，当它存储的数据改变时，依赖于这个 Signal 的计算函数Computer标注“脏” 状态、副作用Effector可以自动运行
// 采用所谓的 “推后拉” 模型：“推” 阶段，在 Signal 变为 “脏”（即其值发生了改变）时，会递归地把 “脏” 状态传递到依赖它的所有计算函数Computer上，所有潜在的重新计算都被推迟，直到显式地 Operate 某个 Computer 的值。

type signal struct {
	value  any
	effect func(any)
	child  []int
}

// Computer 衍生 - 衍生能缓存计算结果，避免重复的计算
// 惰性求值（lazy evaluate）- 只有被使用到的才会计算结果

type computer struct {
	//所有的Signal父代
	parentSignal []int
	//求值链
	evaluateChain []int
	//求值函数
	evaluate func() any
	//Reactions 反应 - 反应是数据更新时的监听器，监视值修改后，立即执行
	effect func(any)
	//false-需计算值，true-无需计算值
	renovate bool
	value    any
}

func (c *computer) eval() {
	c.value = c.evaluate()
	c.renovate = true
	if c.effect != nil {
		c.effect(c.value)
	}
}

type setSignalMsg struct {
	index int
	val   any
}

type operateComputerMsg struct {
	c  int
	do func(any)
}

// Universe 响应式状态管理,放置在一同协程中运算，规避锁的影响
// 适用于生产者慢，消费者快的场景

type Universe struct {
	signalSet           []signal
	computerSet         []computer
	setSignalChan       chan setSignalMsg
	operateComputerChan chan operateComputerMsg
	closeOnce           sync.Once
	// 0: stop 1: live
	stopFlag int32
	stopChan chan struct{}
}

func NewUniverse() *Universe {
	u := &Universe{
		//signalSet[0] 不使用
		signalSet: make([]signal, 1),
		//computerSet[0] 不使用
		computerSet:         make([]computer, 1),
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
		u.signalSet = nil
		u.computerSet = nil
		close(u.setSignalChan)
		close(u.operateComputerChan)
	})
}

// 线程不安全
func (u *Universe) NewSignal(v any, effect func(any)) int {
	u.signalSet = append(u.signalSet, signal{value: v, effect: effect})
	return -len(u.signalSet) + 1
}

// SetSignal 非原子性，乱序执行。
func (u *Universe) SetSignal(n int, a any) {
	if atomic.LoadInt32(&u.stopFlag) != 0 {
		u.setSignalChan <- setSignalMsg{n, a}
	}
}

// GetSignal 线程不安全,需控制在evaluate函数内运行
func (u *Universe) GetSignal(n int) any {
	return u.signalSet[-n].value
}

// GetComputer 线程不安全,需控制在evaluate函数内运行
func (u *Universe) GetComputer(n int) any {
	c := &u.computerSet[n]
	if c.renovate {
		return c.value
	}
	for _, v := range c.evaluateChain {
		k := &u.computerSet[v]
		if !k.renovate {
			k.eval()
		}
	}
	c.eval()
	return c.value
}

// NewComputer 线程不安全
func (u *Universe) NewComputer(effect func(any), warp func(*Universe) (func() any, []int)) int {
	evaluate, p := warp(u)
	c := computer{
		evaluate: evaluate,
		effect:   effect,
		renovate: false,
	}
	j := len(u.computerSet)
	for _, v := range p {
		if v == 0 {
			panic("NewComputer: 不为0")
		}
		// 正负用于判断 signal 或 computer
		if v < 0 {
			ss := &u.signalSet[-v]
			ss.child = append(ss.child, j)
			c.parentSignal = append(c.parentSignal, v)
		} else {
			k := u.computerSet[v]
			c.parentSignal = append(c.parentSignal, k.parentSignal...)
			for _, i := range c.parentSignal {
				ss := &u.signalSet[-i]
				ss.child = append(ss.child, j)
			}
			c.evaluateChain = append(c.evaluateChain, k.evaluateChain...)
			c.evaluateChain = append(c.evaluateChain, v)
		}
	}
	c.parentSignal = RemoveDuplicates(c.parentSignal)
	c.evaluateChain = UniqueWithoutSort(c.evaluateChain)
	u.computerSet = append(u.computerSet, c)
	return j
}

// Operate 非原子性，乱序执行。
func (u *Universe) Operate(n int, do func(any)) {
	if atomic.LoadInt32(&u.stopFlag) != 0 {
		u.operateComputerChan <- operateComputerMsg{c: n, do: do}
	}
}

func (u *Universe) Run() {
	if len(u.signalSet) > 1 {
		for i := range u.signalSet[1:] {
			ss := &u.signalSet[i]
			ss.child = UniqueWithoutSort(ss.child)
		}
	}
	for i := range u.computerSet[1:] {
		u.computerSet[i].parentSignal = nil
	}
	go func() {
		for {
			select {
			case msg := <-u.setSignalChan:
				s := &u.signalSet[-msg.index]
				s.value = msg.val
				for _, v := range s.child {
					u.computerSet[v].renovate = false
				}
				if s.effect != nil {
					s.effect(msg.val)
				}
			case msg := <-u.operateComputerChan:
				if msg.do != nil {
					msg.do(u.GetComputer(msg.c))
				} else {
					u.GetComputer(msg.c)
				}
			case <-u.stopChan:
				return
			}
		}
	}()
}

// https://zhuanlan.zhihu.com/p/691797618
// https://developer.aliyun.com/article/1218778
