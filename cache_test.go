package utils

import (
	"sync"
	"testing"
)

func TestIdempotentCache(t *testing.T) {
	fn := func(d []byte) any {
		return len(d)
	}
	ic := NewIdempotentCache(12, 0x0102030405060708, fn)
	key1 := []byte("127.0.0.1")
	key2 := []byte("192.168.0.1")
	t.Log(ic.Get(key1))
	t.Log(ic.Get(key2))
	t.Log(ic.Get(key1))
}
func BenchmarkMap(b *testing.B) {
	var m sync.Map
	m.Store("127.0.0.1", true)
	for i := 0; i < b.N; i++ {
		m.Load("127.0.0.1")
	}
}

func BenchmarkIdempotentCacheGet(b *testing.B) {
	fn := func(d []byte) any {
		var compute int
		for i := 0; i < 100; i++ {
			compute = len(d) ^ i
		}
		return compute
	}
	ic := NewIdempotentCache(5, 0x0102030405060708, fn)
	key := []byte("127.0.0.1")
	for i := 0; i < b.N; i++ {
		ic.Get(key)
	}
}
