package utils

import (
	"sync"
	"sync/atomic"
)

var defaultPool sync.Pool
var defaultByteSize = uint64(1024)

func init() {
	defaultPool.New = func() interface{} {
		return &[]byte{}
	}
}

func AllocSlice() *[]byte {
	v := defaultPool.Get()
	if v != nil {
		return v.(*[]byte)
	}
	b := make([]byte, 0, atomic.LoadUint64(&defaultByteSize))
	return &b
}
func FreeSlice(x *[]byte) {
	defaultPool.Put(x)
}

func ChangeDefaultByteSize(n uint64) {
	atomic.StoreUint64(&defaultByteSize, n)
}

// https://github.com/valyala/bytebufferpool
