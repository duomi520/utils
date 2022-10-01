package utils

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkMathRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = rand.Int()
	}
}

func BenchmarkRuntimeRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FastRand()
	}
}

func BenchmarkTimeNow(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = time.Now()
	}
}

func BenchmarkRuntimeNanotime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Nanotime()
	}
}
