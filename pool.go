package utils

import (
	"sync"
	"sync/atomic"
)

var defaultPool sync.Pool
var defaultByteSize uint64 = 1024

func GetSlice() *[]byte {
	v := defaultPool.Get()
	if v != nil {
		return v.(*[]byte)
	}
	b := make([]byte, 0, atomic.LoadUint64(&defaultByteSize))
	return &b
}
func PutSlice(x *[]byte) {
	defaultPool.Put(x)
}

// https://github.com/valyala/bytebufferpool
