package utils

import (
	"sync/atomic"
	"time"
)

// RollingWindow 环形滑动窗口统计
type RollingWindow struct {
	//环形窗口总数
	//2^9=512 例：512个窗口，每个窗口约8ms 合计：512*8,388,608=4,294,967,296 约4.3秒
	bucketsCount int
	//滑动统计样本窗口数量
	interval int
	//单个窗口长度，环形窗口总长度
	//2^3=8,2^10=1024 例：1024*1024*8=8,388,608 约8ms  10+10+3=23
	bucketleghthPow, ringPow int
	mask                     int64
	array                    []int64
	round                    []int64
	//Padding
	_ [7]int64
	//0 表示未锁定，非1表示锁定
	spinLock int64
}

// NewRollingWindow
func NewRollingWindow(totalPow, interval, bucketleghthPow int) *RollingWindow {
	r := &RollingWindow{
		bucketsCount:    1 << totalPow,
		interval:        interval,
		bucketleghthPow: bucketleghthPow,
		ringPow:         bucketleghthPow + totalPow,
		mask:            1<<(bucketleghthPow+totalPow) - 1,
	}
	r.array = make([]int64, r.bucketsCount)
	r.round = make([]int64, r.bucketsCount)
	return r
}

// Add
func (r *RollingWindow) Add(n int64) {
	now := time.Now().UnixNano()
	NewRound := now >> r.ringPow
	pos := (int)(now&r.mask) >> r.bucketleghthPow
	if atomic.LoadInt64(&r.round[pos]) == NewRound {
		atomic.AddInt64(&r.array[pos], n)
		return
	}
	var count int
	for {
		if atomic.CompareAndSwapInt64(&r.spinLock, 0, 1) {
			if atomic.LoadInt64(&r.round[pos]) == NewRound {
				atomic.AddInt64(&r.array[pos], n)
			} else {
				atomic.StoreInt64(&r.array[pos], n)
				atomic.StoreInt64(&r.round[pos], NewRound)
			}
			atomic.StoreInt64(&r.spinLock, 0)
			return
		}
		count++
		if count > 10000 {
			panic("RollingWindow.Add 自旋锁死锁")
		}
	}

}

// Store
func (r *RollingWindow) Store(n int64) {
	now := time.Now().UnixNano()
	NewRound := now >> r.ringPow
	pos := (int)(now&r.mask) >> r.bucketleghthPow
	atomic.StoreInt64(&r.array[pos], n)
	atomic.StoreInt64(&r.round[pos], NewRound)
}

func (r *RollingWindow) Sampling() []int {
	now := time.Now().UnixNano()
	offset := now >> r.ringPow
	pre := (int)(now&r.mask)>>r.bucketleghthPow - 1
	var l []int
	if pre >= r.interval {
		for i := pre - r.interval; i < pre; i++ {
			if r.round[i] == offset {
				l = append(l, i)
			}
		}
	} else {
		for i := r.bucketsCount + pre - r.interval; i < r.bucketsCount; i++ {
			if r.round[i] == offset-1 {
				l = append(l, i)
			}
		}
		if pre > 0 {
			for i := range r.round[:pre] {
				if r.round[i] == offset {
					l = append(l, i)
				}
			}
		}
	}
	return l
}

// https://zhuanlan.zhihu.com/p/693443092
// https://www.cnblogs.com/luoxn28/p/11109144.html
// https://www.jianshu.com/p/9cb6aa788520
