package utils

import (
	"strings"
	"testing"
)

var converTestString = strings.Repeat("a", 1024)

func BenchmarkTestString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := []byte(converTestString)
		_ = string(b)
	}
}

func BenchmarkTestBytesToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b := StringToBytes(converTestString)
		_ = BytesToString(b)
	}
}
