package utils

import (
	"runtime"
	"sync/atomic"
)

//LockList 加锁读写数组，适用不频繁写的场景
//修改时新数据原子替换旧数据地址，旧数据由GC回收。
type LockList struct {
	//0-unlock 1-lock
	mutex int64
	slice atomic.Value
}

//NewLockList 新
func NewLockList() *LockList {
	l := LockList{}
	var data []any
	l.slice.Store(data)
	return &l
}

//Add 增加
func (l *LockList) Add(element any) {
	for {
		if atomic.CompareAndSwapInt64(&l.mutex, 0, 1) {
			base := l.slice.Load().([]any)
			size := len(base)
			data := make([]any, size+1)
			copy(data[:size], base)
			data[size] = element
			l.slice.Store(data)
			atomic.StoreInt64(&l.mutex, 0)
			return
		}
		runtime.Gosched()
	}
}

//Remove 移除
func (l *LockList) Remove(judge func(any) bool) {
	for {
		if atomic.CompareAndSwapInt64(&l.mutex, 0, 1) {
			base := l.slice.Load().([]any)
			size := len(base)
			data := make([]any, 0, size)
			for i := 0; i < size; i++ {
				if !judge(base[i]) {
					data = append(data, base[i])
				}
			}
			l.slice.Store(data)
			atomic.StoreInt64(&l.mutex, 0)
			return
		}
		runtime.Gosched()
	}
}

//List 列
func (l *LockList) List() []any {
	return l.slice.Load().([]any)
}

// https://github.com/yireyun/go-queue
// https://github.com/Workiva/go-datastructures
