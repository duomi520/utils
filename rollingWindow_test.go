package utils

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type testBucket struct {
	a []int64
}

func (b *testBucket) New(n int) {
	b.a = make([]int64, n)
}
func (b *testBucket) Add(p int, data []int64) {
	atomic.AddInt64(&b.a[p], data[0])
}
func (b *testBucket) Store(p int, data []int64) {
	atomic.StoreInt64(&b.a[p], data[0])

}
func (b *testBucket) Metric(r *RollingWindow, t time.Duration) func() time.Duration {
	f := func() time.Duration {
		var count int64
		l := r.Sampling()
		for _, v := range l {
			count = count + b.a[v]
		}
		fmt.Println(count, l)
		return t
	}
	return f
}
func TestRollingWindowAdd(t *testing.T) {
	b := &testBucket{}
	//2^4=16 ,2^27=134,217,728 约134 ms
	r := NewRollingWindow(4, 6, 27, b)
	for range 10 {
		time.Sleep(134 * time.Millisecond)
		r.Add([]int64{1})
	}
	fmt.Println(b.a)
	fmt.Println(r.round)
	time.Sleep(500 * time.Millisecond)
	for range 80 {
		time.Sleep(13 * time.Millisecond)
		r.Add([]int64{1})
	}
	fmt.Println(b.a)
	fmt.Println(r.round)
}

// [1 0 0 0 0 0 0 1 1 1 1 1 1 1 1 1]
// [811149755 0 0 0 0 0 0 811149754 811149754 811149754 811149754 811149754 811149754 811149754 811149754 811149754]
// [1 0 0 1 10 10 10 10 10 10 10 9 1 1 1 1]
// [811149755 0 0 811149755 811149755 811149755 811149755 811149755 811149755 811149755 811149755 811149755 811149754 811149754 811149754 811149754]
func TestRollingWindowSampling(t *testing.T) {
	b := &testBucket{}
	r := NewRollingWindow(4, 6, 27, b)
	for range 10 {
		time.Sleep(134 * time.Millisecond)
		r.Add([]int64{1})
	}
	time.Sleep(134 * time.Millisecond)
	b.Metric(r, time.Millisecond)
	fmt.Println(b.a)
	fmt.Println(r.round)
	time.Sleep(500 * time.Millisecond)
	for range 80 {
		time.Sleep(13 * time.Millisecond)
		r.Add([]int64{1})
	}
	pre := (int)(time.Now().UnixNano()&r.mask)>>r.bucketleghthPow - 2
	if pre > -1 {
		r.round[pre] = 0
		r.round[pre+1] = 0
	} else {
		r.round[r.bucketsCount-1] = 0
	}
	f := b.Metric(r, time.Millisecond)
	f()
	fmt.Println(pre, b.a)
	fmt.Println(r.round)
}

/*
7 [9 10 11 12 13 14 15]
[0 0 0 0 0 0 1 1 1 1 1 1 1 1 1 1]
[0 0 0 0 0 0 811149847 811149847 811149847 811149847 811149847 811149847 811149847 811149847 811149847 811149847]
50 [4 5 6 7 8]
9 [0 0 0 0 10 10 10 10 10 10 10 10 1 1 1 1]
[0 0 0 0 811149848 811149848 811149848 811149848 811149848 0 0 811149848 811149847 811149847 811149847 811149847]
*/
