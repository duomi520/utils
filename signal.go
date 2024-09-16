package utils

import (
	"fmt"
	"slices"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

var globalIncreases int64

type Parent interface {
	Load() any
}

// Signal 信号-状态单向传播
// 一个数据的容器，当它存储的数据改变时，依赖于这个 Signal 的计算函数Computer标注“脏” 状态、副作用Effector可以自动运行
// 采用所谓的 “推后拉” 模型：“推” 阶段，在 Signal 变为 “脏”（即其值发生了改变）时，会递归地把 “脏” 状态传递到依赖它的所有计算函数Computer上，所有潜在的重新计算都被推迟，直到显式地 Operate 某个 Computer 的值。

type Signal struct {
	//去重标记
	id int64
	atomic.Value
	//所有的后代
	child []*Computer
	//Reactions 反应 - 反应是数据更新时的监听器，监视值修改后，立即执行
	effect func(any)
}

// 线程不安全
func NewSignal(v any, do func(any)) *Signal {
	s := &Signal{
		id:     atomic.AddInt64(&globalIncreases, 1),
		effect: do,
	}
	s.Store(v)
	return s
}

// Set 非原子性，有一定的延迟。
func (s *Signal) Set(u *Universe, a any) {
	if atomic.LoadInt32(&u.stopFlag) != 0 {
		u.SetSignalChan <- setSignalMsg{s, a}
	}
}

// Computer 衍生 - 衍生能缓存计算结果，避免重复的计算
// 惰性求值（lazy evaluate）- 只有被使用到的才会计算结果
type Computer struct {
	//去重标记
	id int64
	//计算层
	tier int32
	//所有的Signal父代
	allParentSignal []*Signal
	//所有的Computer父代
	allParentComputer []*Computer
	//false-需计算值，true-无需计算值
	renovate bool
	value    any
	//上一个父代
	parent []Parent
	//求值函数
	evaluate func(...Parent) any
	//Reactions 反应 - 反应是数据更新时的监听器，监视值修改后，立即执行
	effect func(any)
}

type ByTier []*Computer

func (t ByTier) Len() int           { return len(t) }
func (t ByTier) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t ByTier) Less(i, j int) bool { return t[i].tier < t[j].tier }

// NewComputer 线程不安全
func NewComputer(effect func(any), evaluate func(...Parent) any, p ...Parent) *Computer {
	c := &Computer{
		id:       atomic.AddInt64(&globalIncreases, 1),
		renovate: false,
		parent:   p,
		evaluate: evaluate,
		effect:   effect,
	}
	for _, v := range p {
		switch p := v.(type) {
		case *Signal:
			p.child = append(p.child, c)
			c.allParentSignal = append(c.allParentSignal, p)
		case *Computer:
			c.allParentComputer = append(c.allParentComputer, p.allParentComputer...)
			c.allParentSignal = append(c.allParentSignal, p.allParentSignal...)
		default:
			panic(fmt.Sprintf("NewComputer: 无效的参数类型:%v", v))
		}

	}
	//allParentSignal去重
	if len(c.allParentSignal) > 0 {
		slices.SortFunc(c.allParentSignal, func(a, b *Signal) int {
			return int(a.id - b.id)
		})
		var ps []*Signal
		for i := 1; i < len(c.allParentSignal); i++ {
			if c.allParentSignal[i].id != c.allParentSignal[i-1].id {
				ps = append(ps, c.allParentSignal[i])
			}
		}
		c.allParentSignal = ps
	}
	//allParentComputer去重
	if len(c.allParentComputer) > 0 {
		//allParentComputer去重
		slices.SortFunc(c.allParentComputer, func(a, b *Computer) int {
			return int(a.id - b.id)
		})
		var pc []*Computer
		for i := 1; i < len(c.allParentComputer); i++ {
			if c.allParentComputer[i].id != c.allParentComputer[i-1].id {
				pc = append(pc, c.allParentComputer[i])
			}
		}
		c.allParentComputer = pc
		//allParentComputer排序
		sort.Sort(ByTier(c.allParentComputer))
		//设置计算层
		c.tier = c.allParentComputer[len(c.allParentComputer)-1].tier + 1
	}
	return c
}

// Load 线程不安全,需控制在evaluate函数内运行
func (c *Computer) Load() any {
	if c.renovate {
		return c.value
	}
	for _, v := range c.allParentComputer {
		if !v.renovate {
			v.evaluate(v.parent...)
		}
	}
	c.value = c.evaluate(c.parent...)
	return c.value
}

// Operate 非原子性，有一定的延迟。
func (c *Computer) Operate(u *Universe) {
	if atomic.LoadInt32(&u.stopFlag) != 0 {
		u.ComputerChan <- c
	}
}

type setSignalMsg struct {
	signal *Signal
	val    any
}

// Universe 响应式状态管理,放置在一同协程中运算，规避锁的影响
// 适用于生产者慢，消费者快的场景

type Universe struct {
	SetSignalChan chan setSignalMsg
	ComputerChan  chan *Computer
	closeOnce     sync.Once
	// 0: stop 1: live
	stopFlag int32
	stopChan chan struct{}
}

func NewUniverse() *Universe {
	u := &Universe{
		SetSignalChan: make(chan setSignalMsg, 128),
		ComputerChan:  make(chan *Computer, 128),
		stopFlag:      1,
		stopChan:      make(chan struct{}),
	}
	return u
}

func (u *Universe) Run() {
	go func() {
		for {
			select {
			case msg := <-u.SetSignalChan:
				msg.signal.Store(msg.val)
				for _, v := range msg.signal.child {
					v.renovate = false
				}
				if msg.signal.effect != nil {
					msg.signal.effect(msg.signal.Load())
				}
			case c := <-u.ComputerChan:
				if c.renovate {
					if c.effect != nil {
						c.effect(c.value)
					}
				} else {
					if c.effect != nil {
						c.effect(c.Load())
					}
				}
			case <-u.stopChan:
				return
			}
		}
	}()
}

// Close 关闭
func (u *Universe) Close() {
	u.closeOnce.Do(func() {
		atomic.StoreInt32(&u.stopFlag, 0)
		close(u.stopChan)
		time.Sleep(3 * time.Second)
		close(u.SetSignalChan)
		close(u.ComputerChan)
	})
}

// https://zhuanlan.zhihu.com/p/691797618
// https://developer.aliyun.com/article/1218778
