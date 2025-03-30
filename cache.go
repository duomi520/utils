package utils

import (
	"sync/atomic"
	"unsafe"
)

type Cell struct {
	wyhash  uint64
	fnv1a   uint64
	timeout int64
	out     any
}

// IdempotentCache 幂等函数缓存，幂等方法，是指可以使用相同参数重复执行，并能获得相同结果的函数
type IdempotentCache[T string | []byte] struct {
	size uint64
	//用于hash的种子
	seed uint64
	buf  []unsafe.Pointer
	do   func(T) any
}

// Init 初始化 power表示缓存最大数量的的2的次方，seed表示hash的种子，keep 表示缓存保留时间，cycle 表示周期清理任务时间，do表示要缓存的幂等函数
func (ic *IdempotentCache[T]) Init(power, seed uint64, do func(T) any) {
	//缓存的大小，使用2的power次方作为大小。
	ic.size = 2 ^ power
	ic.seed = seed
	ic.buf = make([]unsafe.Pointer, ic.size)
	ic.do = do
}

// Get 用于获取缓存中的结果
func (ic *IdempotentCache[T]) Get(in T) any {
	h := Hash64WY(in, ic.seed)
	f := Hash64FNV1A(in)
	//取余
	index := h & (ic.size - 1)
	v := atomic.LoadPointer(&ic.buf[index])
	if v != nil {
		cell := (*Cell)(v)
		if cell.wyhash == h && cell.fnv1a == f {
			return cell.out
		}
	}
	c := &Cell{wyhash: h, fnv1a: f, out: ic.do(in)}
	atomic.StorePointer(&ic.buf[index], unsafe.Pointer(c))
	return c.out
}

// Remove 移除缓存
func (ic *IdempotentCache[T]) Remove(in T) {
	h := Hash64WY(in, ic.seed)
	f := Hash64FNV1A(in)
	//取余
	index := h & (ic.size - 1)
	v := atomic.SwapPointer(&ic.buf[index], nil)
	if v != nil {
		cell := (*Cell)(v)
		if cell.wyhash == h && cell.fnv1a == f {
			return
		}
		atomic.StorePointer(&ic.buf[index], v)
	}
}

// https://github.com/cespare/xxhash
// https://zhuanlan.zhihu.com/p/624248354
// https://zhuanlan.zhihu.com/p/466139082
