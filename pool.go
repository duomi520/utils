package utils

import (
	"bytes"
	"sync"
	"sync/atomic"
)

const (
	poolSteps = 26

	// 64大小的数据可以一次性被加载到 CPU 缓存行中。
	minBitSize int64 = 64

	calibrateCallsThreshold int64 = 42000
)

// Pool
type Pool struct {
	slicePool, byteBufferPool sync.Pool
	array                     [poolSteps]int64
	//Padding
	_           [7]int64
	calibrating int64
	//Padding
	_           [7]int64
	defaultSize int64
}

func (p *Pool) AllocSlice() *[]byte {
	v := p.slicePool.Get()
	if v != nil {
		return v.(*[]byte)
	}
	b := make([]byte, 0, max(atomic.LoadInt64(&p.defaultSize), minBitSize))
	return &b
}
func (p *Pool) FreeSlice(x *[]byte) {
	idx := Log2Up(uint32(len(*x)))
	if atomic.AddInt64(&p.array[idx], 1) > calibrateCallsThreshold {
		p.calibrate()
	}
	// 重置切片长度为 0
	*x = (*x)[:0]
	n := atomic.LoadInt64(&p.defaultSize)
	if n == 0 || cap(*x) <= int(n) {
		p.slicePool.Put(x)
	}
}

func (p *Pool) AllocBuffer() *bytes.Buffer {
	v := p.byteBufferPool.Get()
	if v != nil {
		return v.(*bytes.Buffer)
	}
	b := make([]byte, 0, max(atomic.LoadInt64(&p.defaultSize), minBitSize))
	return bytes.NewBuffer(b)
}

func (p *Pool) FreeBuffer(x *bytes.Buffer) {
	idx := Log2Up(uint32(x.Len()))
	if atomic.AddInt64(&p.array[idx], 1) > calibrateCallsThreshold {
		p.calibrate()
	}
	// 重置 Buffer
	x.Reset()
	n := atomic.LoadInt64(&p.defaultSize)
	if n == 0 || x.Cap() <= int(n) {
		p.byteBufferPool.Put(x)
	}
}

// 校准操作
func (p *Pool) calibrate() {
	// 避免并发
	if !atomic.CompareAndSwapInt64(&p.calibrating, 0, 1) {
		return
	}
	var list [poolSteps]int64
	var total, sum int64
	// 读出总数
	for i := range poolSteps {
		list[i] = atomic.SwapInt64(&p.array[i], 0)
		total += list[i]
	}
	total = total * 9 / 10
	// 统计
	var pos int
	for i := range poolSteps {
		if sum > total {
			break
		}
		sum += list[i]
		pos = i
	}
	// fmt.Println(pos, list)
	// 保存对应值
	atomic.StoreInt64(&p.defaultSize, max(1<<pos, minBitSize))
	atomic.StoreInt64(&p.calibrating, 0)
}

// Log2Up 取上一个2的对数
func Log2Up(x uint32) int {
	var ans int
	if (x & 0xffff0000) > 0 {
		ans += 16
		x &= 0xffff0000
	}
	if (x & 0xff00ff00) > 0 {
		ans += 8
		x &= 0xff00ff00
	}
	if (x & 0xf0f0f0f0) > 0 {
		ans += 4
		x &= 0xf0f0f0f0
	}
	if (x & 0xcccccccc) > 0 {
		ans += 2
		x &= 0xcccccccc
	}
	if (x & 0xaaaaaaaa) > 0 {
		ans += 1
		x &= 0xaaaaaaaa
	}
	return ans
}

// https://github.com/valyala/bytebufferpool
// https://github.com/oxtoacart/bpool
// https://zhuanlan.zhihu.com/p/370848384
