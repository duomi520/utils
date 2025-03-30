package utils

import (
	"errors"
	"sync/atomic"
	"time"
)

// Limiter 限流器 Token Bucket(令牌桶)
// 每隔一段时间加入一批令牌，达到上限后，不再增加。
type TokenBucketLimiter struct {
	//限流器速率，每秒处理的令牌数
	LimitRate int64
	//限流器大小，存放令牌的最大值
	LimitSize int64
	//加入的时间间隔
	Snippet time.Duration
	//Padding
	_ [7]int64
	//令牌
	tokens int64
	//退出标志 1-退出
	stopFlag int32
}

// NewTokenBucketLimiter limitRate, limitSize,snippet数值较小时，准确度低。
func NewTokenBucketLimiter(limitRate, limitSize int64, snippet time.Duration) *TokenBucketLimiter {
	t := &TokenBucketLimiter{
		LimitRate: limitRate,
		LimitSize: limitSize,
		Snippet:   snippet,
		tokens:    limitSize,
		stopFlag:  0,
	}
	if t.Snippet == 0 {
		//默认0.1秒
		t.Snippet = 100 * time.Millisecond
	}
	return t
}

// Wait 等待,申请n个令牌，取不到足够数量时返回错误。
func (t *TokenBucketLimiter) Take(n int64) error {
	if atomic.LoadInt32(&t.stopFlag) == 1 {
		return nil
	}
	new := atomic.AddInt64(&t.tokens, -n)
	if new > -1 {
		return nil
	}
	atomic.AddInt64(&t.tokens, n)
	return errors.New("rate limit")
}

func (t *TokenBucketLimiter) Run() {
	for {
		if atomic.LoadInt32(&t.stopFlag) == 1 {
			return
		}
		time.Sleep(t.Snippet)
		//竟态下，牺牲准确度。
		new := atomic.AddInt64(&t.tokens, t.LimitRate)
		if new > t.LimitSize {
			atomic.StoreInt64(&t.tokens, t.LimitSize)
		}
		if new < t.LimitRate {
			atomic.StoreInt64(&t.tokens, t.LimitRate)
		}
	}
}

func (t *TokenBucketLimiter) Task() time.Duration {
	if atomic.LoadInt32(&t.stopFlag) == 1 {
		return 0
	}
	//竟态下，牺牲准确度。
	new := atomic.AddInt64(&t.tokens, t.LimitRate)
	if new > t.LimitSize {
		atomic.StoreInt64(&t.tokens, t.LimitSize)
	}
	if new < t.LimitRate {
		atomic.StoreInt64(&t.tokens, t.LimitRate)
	}
	return t.Snippet
}

// Close 关闭。
func (t *TokenBucketLimiter) Close() {
	atomic.StoreInt32(&t.stopFlag, 1)
}

// https://golang.org/x/time/rate
// https://studygolang.com/articles/27454#reply0
// https://zhuanlan.zhihu.com/p/89820414
// https://github.com/juju/ratelimit
// https://github.com/uber-go/ratelimit
// https://mp.weixin.qq.com/s/T_LvVfAOzgANO1XSCViJrw
// https://mp.weixin.qq.com/s/YCvUTwpe0jUdwcKyQQj7hA
