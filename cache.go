package utils

import (
	"sync/atomic"
	"time"
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
	keep time.Duration
	buf  []unsafe.Pointer
	do   func(T) any
}

// Init 初始化 power表示缓存最大数量的的2的次方，seed表示hash的种子，keep 表示缓存保留时间，cycle 表示周期清理任务时间，do表示要缓存的幂等函数
func (ic *IdempotentCache[T]) Init(power, seed uint64, keep, cycle time.Duration, t *Timing, do func(T) any) {
	//缓存的大小，使用2的power次方作为大小。
	ic.size = 2 ^ power
	ic.seed = seed
	ic.keep = keep
	ic.buf = make([]unsafe.Pointer, ic.size)
	ic.do = do
	if t != nil {
		t.AddTask(
			time.Now().Add(cycle),
			func() time.Duration {
				return func() time.Duration {
					// 费资源，BitMap是否会更快？
					now := time.Now().Unix()
					for i := range ic.buf {
						old := atomic.LoadPointer(&ic.buf[i])
						if old != nil {
							cell := (*Cell)(old)
							over := atomic.LoadInt64(&cell.timeout)
							if now < over {
								atomic.CompareAndSwapPointer(&ic.buf[i], old, nil)
							}
						}
					}
					return cycle
				}()
			},
		)
	}
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
			//存在竟态，有可能会几个协程同时写入时间，不影响使用。
			atomic.StoreInt64(&cell.timeout, time.Now().Add(ic.keep).Unix())
			return cell.out
		}
	}
	c := &Cell{wyhash: h, fnv1a: f, timeout: time.Now().Add(ic.keep).Unix(), out: ic.do(in)}
	atomic.StorePointer(&ic.buf[index], unsafe.Pointer(c))
	return c.out
}

// https://github.com/cespare/xxhash
// https://zhuanlan.zhihu.com/p/624248354
// https://zhuanlan.zhihu.com/p/466139082
