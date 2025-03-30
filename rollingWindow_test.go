package utils

import (
	"fmt"
	"testing"
	"time"
)

func testMetric( r *RollingWindow, t time.Duration) func() time.Duration {
	return func() time.Duration {
		var count int64
		l := r.Sampling()
		for _, v := range l {
			count = count + r.array[v]
		}
		fmt.Println(count, l)
		return t
	}
}
func TestRollingWindowStore(t *testing.T) {
	//2^4=16 ,2^27=134,217,728 约134 ms
	r := NewRollingWindow(4, 6, 27)
	for i := range 10 {
		time.Sleep(134 * time.Millisecond)
		r.Store(int64(i))
	}
	time.Sleep(2 * time.Second)
	for i := range 5 {
		time.Sleep(134 * time.Millisecond)
		r.Store(int64(i + 100))
	}
	fmt.Println(r.array)
	fmt.Println(r.round)
}

/*
[2 3 4 5 6 7 8 100 101 102 103 104 0 0 0 1]
[811638564 811638564 811638564 811638564 811638564 811638564 811638564 811638565 811638565 811638565 811638565 811638565 0 0 811638563 811638563]
*/

func TestRollingWindowAdd(t *testing.T) {
	//2^4=16 ,2^27=134,217,728 约134 ms
	r := NewRollingWindow(4, 6, 27)
	for range 10 {
		time.Sleep(134 * time.Millisecond)
		r.Add(1)
	}
	fmt.Println(r.array)
	fmt.Println(r.round)
	time.Sleep(500 * time.Millisecond)
	for range 80 {
		time.Sleep(13 * time.Millisecond)
		r.Add(1)
	}
	fmt.Println(r.array)
	fmt.Println(r.round)
}

/*
[0 0 0 0 0 1 1 1 1 1 1 1 1 1 1 0]
[0 0 0 0 0 811597425 811597425 811597425 811597425 811597425 811597425 811597425 811597425 811597425 811597425 0]
[0 2 10 10 10 10 10 10 10 8 1 1 1 1 1 0]
[0 811597426 811597426 811597426 811597426 811597426 811597426 811597426 811597426 811597426 811597425 811597425 811597425 811597425 811597425 0]
*/
func TestRollingWindowSampling(t *testing.T) {
	r := NewRollingWindow(4, 6, 27)
	for range 10 {
		time.Sleep(134 * time.Millisecond)
		r.Add(1)
	}
	time.Sleep(134 * time.Millisecond)
	testMetric(r, time.Millisecond)
	fmt.Println(r.array)
	fmt.Println(r.round)
	time.Sleep(500 * time.Millisecond)
	for range 80 {
		time.Sleep(13 * time.Millisecond)
		r.Add(1)
	}
	pre := (int)(time.Now().UnixNano()&r.mask)>>r.bucketleghthPow - 2
	if pre > -1 {
		r.round[pre] = 0
		r.round[pre+1] = 0
	} else {
		r.round[r.bucketsCount-1] = 0
	}
	f := testMetric(r, time.Millisecond)
	f()
	fmt.Println(pre, r.array)
	fmt.Println(r.round)
}

/*
[1 1 0 0 0 0 0 0 1 1 1 1 1 1 1 1]
[811597461 811597461 0 0 0 0 0 0 811597460 811597460 811597460 811597460 811597460 811597460 811597460 811597460]
50 [6 7 8 9 10]
11 [1 1 0 0 0 2 10 10 10 10 10 10 10 8 1 1]
[811597461 811597461 0 0 0 811597461 811597461 811597461 811597461 811597461 811597461 0 0 811597461 811597460 811597460]
*/


func BenchmarkRollingWindowStore(b *testing.B) {
	r := NewRollingWindow(4, 6, 27)
	for i := 0; i < b.N; i++ {
		r.Store(1)
	}
}

func BenchmarkRollingWindowAdd(b *testing.B) {
	r := NewRollingWindow(4, 6, 27)
	for i := 0; i < b.N; i++ {
		r.Add(1)
	}
}