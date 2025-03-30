package utils

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestLog2Up(t *testing.T) {
	var tests = []struct {
		arg    uint32
		result int
	}{
		{0, 0}, {1, 0}, {66, 6}, {158, 7}, {1026, 10}, {2555, 11},
	}
	for i := range tests {
		if Log2Up(tests[i].arg) != tests[i].result {
			t.Fatal(i, tests[i].arg, tests[i].result)
		}
	}
}

func TestPool(t *testing.T) {
	p := &Pool{}
	for i := range 10000 {
		b := p.AllocSlice()
		buf := make([]byte, i)
		*b = append(*b, buf...)
		p.FreeSlice(b)
	}
	fmt.Println(p)
}

// &{{{} 0xc00010f508 12 0xc00010ee08 12 <nil>} {{} <nil> 0 <nil> 0 <nil>} [2 2 4 8 16 32 64 128 256 512 1024 2048 4096 1808 0 0 0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0] 0 [0 0 0 0 0 0 0] 0}
func TestCalibrate(t *testing.T) {
	p := &Pool{}
	for range 1000000 {
		b := p.AllocSlice()
		buf := make([]byte, rand.Intn(10000))
		*b = append(*b, buf...)
		p.FreeSlice(b)
	}
	fmt.Println(p)
}

// &{{{} 0xc00007dc08 12 0xc00007d508 0 <nil>} {{} <nil> 0 <nil> 0 <nil>} [18 16 34 55 123 235 503 953 1995 3905 7876 15933 31743 14060 0 0 0 0 0 0 0 0 0 0 0 0] [0 0 0 0 0 0 0] 0 [0 0 0 0 0 0 0] 8192}
