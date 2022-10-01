package utils

import (
	"strconv"
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
func BenchmarkTestItoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(456789)
	}
}
func BenchmarkTestUItoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		UItoa(uint(456789))
	}
}

func TestIntegerToString(t *testing.T) {
	var tests = []struct {
		in  uint
		out string
	}{
		{0, "0"},
		{123456789, "123456789"},
	}
	for i := range tests {
		s := UItoa(tests[i].in)
		if !strings.EqualFold(s, tests[i].out) {
			t.Errorf("expected %s got %s", tests[i].out, s)
		}
		t.Log(s, tests[i].out)
	}
}
