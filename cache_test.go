package utils

import (
	"sync"
	"testing"
)

func TestIdempotentCache(t *testing.T) {
	fn := func(d []byte) any {
		return len(d)
	}
	ic := &IdempotentCache[[]byte]{}
	ic.Init(12, 0x0102030405060708, fn)
	key1 := []byte("127.0.0.1")
	key2 := []byte("192.168.0.1")
	t.Log(ic.Get(key1))
	t.Log(ic.Get(key2))
	t.Log(ic.Get(key1))
}
func BenchmarkSyncMap(b *testing.B) {
	var m sync.Map
	m.Store("127.0.0.1", true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Load("127.0.0.1")
	}
}

func BenchmarkIdempotentCacheGet(b *testing.B) {
	fn := func(s string) any {
		var compute int
		for i := range 100 {
			compute = len(s) ^ i
		}
		return compute
	}
	ic := &IdempotentCache[string]{}
	ic.Init(5, 0x0102030405060708, fn)
	ic.Get("127.0.0.1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ic.Get("127.0.0.1")
	}
}
