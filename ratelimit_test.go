package utils

import (
	"fmt"
	"testing"
	"time"
)

func BenchmarkWait(b *testing.B) {
	t := NewTokenBucketLimiter(1000, 16*1024, 10*time.Millisecond)
	defer t.Close()
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < b.N; i++ {
		t.Take(1)
	}
}

func TestAllowed(t *testing.T) {
	limiter := NewTokenBucketLimiter(100, 200, 100*time.Millisecond)
	defer limiter.Close()
	go limiter.Run()
	time.Sleep(10 * time.Millisecond)
	count := 0
	fmt.Println(count, limiter)
	for i := 0; i < 500; i++ {
		if err := limiter.Take(1); err == nil {
			count++
		}
	}
	fmt.Println(count, limiter)
	time.Sleep(100 * time.Millisecond)
	fmt.Println(count, limiter)
	for i := 0; i < 500; i++ {
		if err := limiter.Take(1); err == nil {
			count++
		}
	}
	fmt.Println(count, limiter)
}

//0 &{100 200 100ms 200 0}
//200 &{100 200 100ms 0 0}
//200 &{100 200 100ms 100 0}
//300 &{100 200 100ms 0 0}

func TestTake(t *testing.T) {
	limiter := NewTokenBucketLimiter(2, 4, 10*time.Millisecond)
	defer limiter.Close()
	go limiter.Run()
	time.Sleep(10 * time.Millisecond)
	for i := 0; i < 10; i++ {
		prev := time.Now()
		if err := limiter.Take(3); err != nil {
			fmt.Println(i, time.Since(prev), limiter.tokens, err.Error())
		} else {
			fmt.Println(i, time.Since(prev), limiter.tokens)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

/*
0 0s 1
1 83.3µs 1 rate limit
2 0s 1
3 0s 3 rate limit
4 0s 0
5 0s 2 rate limit
6 0s 1
7 0s 0
8 0s 2 rate limit
9 0s 1
*/
