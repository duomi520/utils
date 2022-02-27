package utils

import (
	"log"
	"testing"
)

type lockListTemp struct {
	i int
}

func (l lockListTemp) equal(obj any) bool {
	return l.i == obj.(lockListTemp).i
}
func TestLockList(t *testing.T) {
	l := NewLockList()
	for i := 0; i < 10; i++ {
		l.Add(lockListTemp{i})
	}
	var tests = []struct {
		arg    int
		result []int
	}{
		{5, []int{0, 1, 2, 3, 4, 6, 7, 8, 9}},
		{2, []int{0, 1, 3, 4, 6, 7, 8, 9}},
		{9, []int{0, 1, 3, 4, 6, 7, 8}},
		{0, []int{1, 3, 4, 6, 7, 8}},
		{10, []int{1, 3, 4, 6, 7, 8}},
	}
	for _, v := range tests {
		l.Remove(lockListTemp{v.arg}.equal)
		temp := l.List()
		if len(v.result) != len(l.List()) {
			log.Fatalln(v, temp)
		}
		for i := range temp {
			if v.result[i] != temp[i].(lockListTemp).i {
				log.Fatalln(i, v, temp)
			}
		}
	}
}
func BenchmarkList(b *testing.B) {
	l := NewLockList()
	for i := 0; i < 10; i++ {
		l.Add(lockListTemp{i})
	}
	for i := 0; i < b.N; i++ {
		_ = l.List()
	}
}
