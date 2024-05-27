package utils

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
	"time"
)

type task struct {
	// 下次执行时间,元素在队列中的优先级
	next time.Time
	// 元素在堆中的索引
	index int
	//返回下次执行间隔时间,0 退出
	do func() time.Duration
}

type priorityQueue []task

func (p priorityQueue) Len() int           { return len(p) }
func (p priorityQueue) Less(i, j int) bool { return p[i].next.Before(p[j].next) }
func (p priorityQueue) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p *priorityQueue) Push(x any) {
	n := len(*p)
	item := x.(task)
	item.index = n
	*p = append(*p, item)
}
func (p *priorityQueue) Pop() any {
	old := *p
	n := len(old)
	item := old[n-1]
	// 为了安全性考虑而做的设置
	item.index = -1
	*p = old[0 : n-1]
	return item
}

type Timing struct {
	panicHandler func(error)
	queue        priorityQueue
	AddTaskChan  chan task
	closeOnce    sync.Once
	stopChan     chan struct{}
}

// NewTiming 新建
func NewTiming(p func(error)) *Timing {
	var t = Timing{
		panicHandler: p,
		queue:        make(priorityQueue, 0, 128),
		AddTaskChan:  make(chan task, 128),
		stopChan:     make(chan struct{}),
	}
	go t.run()
	return &t
}

// Stop 停止
func (t *Timing) Stop() {
	t.closeOnce.Do(func() {
		close(t.stopChan)
	})

}

// AddTask 加入任务
func (t *Timing) AddTask(next time.Time, f func() time.Duration) error {
	warp := func(base func() time.Duration) func() time.Duration {
		defer func() {
			buf, r := FormatRecover()
			if r != nil && t.panicHandler != nil {
				t.panicHandler(fmt.Errorf("异常拦截： %s \n%s", r, buf))
			}
		}()
		return base
	}
	select {
	case <-t.stopChan:
		return errors.New("Timing closed")
	case t.AddTaskChan <- task{next: next, do: warp(f)}:
	default:
		return errors.New("join failed")
	}
	return nil
}

func (t *Timing) run() {
	timer := time.NewTimer(time.Second)
	defer timer.Stop()
	heap.Init(&t.queue)
	var interval time.Duration = 64 * 365 * 24 * time.Hour
	for {
		select {
		case v := <-t.AddTaskChan:
			heap.Push(&t.queue, v)
			new := time.Until(v.next)
			if new < interval {
				interval = new
				timer.Reset(interval)
			}
		case <-timer.C:
			if len(t.queue) > 0 {
				v1 := heap.Pop(&t.queue).(task)
				space := v1.do()
				if space > 0 {
					v1.next = time.Now().Add(space)
					heap.Push(&t.queue, v1)
				}
				if len(t.queue) > 0 {
					v2 := heap.Pop(&t.queue).(task)
					interval = time.Until(v2.next)
					heap.Push(&t.queue, v2)
					timer.Reset(interval)
				} else {
					interval = 64 * 365 * 24 * time.Hour
				}
			}
		case <-t.stopChan:
			return
		}
	}
}

// https://www.cnblogs.com/dream397/p/15021120.html
