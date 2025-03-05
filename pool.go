package utils

import (
	"bytes"
	"sync"
	"sync/atomic"
)

var defaultSlicePool sync.Pool
var defaultByteSize = uint64(1024)
var defaultByteBufferPool sync.Pool

func AllocSlice() *[]byte {
	v := defaultSlicePool.Get()
	if v != nil {
		return v.(*[]byte)
	}
	b := make([]byte, 0, atomic.LoadUint64(&defaultByteSize))
	return &b
}
func FreeSlice(x *[]byte) {
	// 重置切片长度为 0
	*x = (*x)[:0]
	defaultSlicePool.Put(x)
}

func AllocBuffer() *bytes.Buffer {
	v := defaultByteBufferPool.Get()
	if v != nil {
		buf := v.(*bytes.Buffer)
		// 重置 Buffer
		buf.Reset()
		return buf
	}
	b := make([]byte, 0, atomic.LoadUint64(&defaultByteSize))
	return bytes.NewBuffer(b)
}

func FreeBuffer(x *bytes.Buffer) {
	// 重置 Buffer
	x.Reset()
	defaultByteBufferPool.Put(x)
}

func ChangeDefaultByteSize(n uint64) {
	atomic.StoreUint64(&defaultByteSize, n)
}

// https://github.com/valyala/bytebufferpool
// https://github.com/oxtoacart/bpool
