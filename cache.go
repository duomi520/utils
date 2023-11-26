package utils

import (
	"bytes"
	"sync/atomic"
)

type Cell struct {
	in  []byte
	out any
}

// IdempotentCache 幂等函数缓存，幂等方法，是指可以使用相同参数重复执行，并能获得相同结果的函数
type IdempotentCache struct {
	power uint64
	size  uint64
	seed  uint64
	buf   []atomic.Value
	do    func([]byte) any
}

// NewIdempotentCache 新建
func NewIdempotentCache(power, seed uint64, do func([]byte) any) *IdempotentCache {
	return &IdempotentCache{
		power: power,
		size:  2 ^ power,
		seed:  seed,
		buf:   make([]atomic.Value, 2^power),
		do:    do,
	}
}

// Get 取
func (ic *IdempotentCache) Get(in []byte) any {
	got := Hash64(in, ic.seed)
	//取余
	index := got & (ic.size - 1)
	v := ic.buf[index].Load()
	if v != nil {
		cell := v.(Cell)
		if bytes.EqualFold(in, cell.in) {
			return cell.out
		}
	}
	var c Cell
	c.in = in
	c.out = ic.do(in)
	ic.buf[index].Store(c)
	return c.out
}
