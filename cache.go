package utils

import (
	"sync/atomic"
)

type Cell struct {
	wyhash uint64
	fnv1a  uint64
	out    any
}

// IdempotentCache 幂等函数缓存，幂等方法，是指可以使用相同参数重复执行，并能获得相同结果的函数
type IdempotentCache[T string | []byte] struct {
	power uint64
	size  uint64
	seed  uint64
	buf   []atomic.Value
	do    func(T) any
}

// Init 初始化
func (ic *IdempotentCache[T]) Init(power, seed uint64, do func(T) any) {
	ic.power = power
	ic.size = 2 ^ power
	ic.seed = seed
	ic.buf = make([]atomic.Value, 2^power)
	ic.do = do
}

// Get 取
func (ic *IdempotentCache[T]) Get(in T) any {
	h := Hash64WY(in, ic.seed)
	f := Hash64FNV1A(in)
	//取余
	index := h & (ic.size - 1)
	v := ic.buf[index].Load()
	if v != nil {
		cell := v.(Cell)
		if cell.wyhash == h && cell.fnv1a == f {
			return cell.out
		}
	}
	var c Cell
	c.wyhash = h
	c.fnv1a = f
	c.out = ic.do(in)
	ic.buf[index].Store(c)
	return c.out
}

// https://github.com/cespare/xxhash
