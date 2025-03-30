package utils

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	var tests = []struct {
		arg    []int
		result []int
	}{
		{[]int{8, 2, 7, 3, 5, 3, 4, 5}, []int{2, 3, 4, 5, 7, 8}},
		{[]int{0, 1, 3, 4, 6, 7, 8, 0}, []int{0, 1, 3, 4, 6, 7, 8}},
	}
	for i := range tests {
		r := RemoveDuplicates(tests[i].arg)
		if !reflect.DeepEqual(r, tests[i].result) {
			t.Fatal(i, r, tests[i].result)
		}
	}
}

func TestUniqueWithoutSort(t *testing.T) {
	var tests = []struct {
		arg    []int
		result []int
	}{
		{[]int{8, 2, 7, 3, 5, 3, 4, 5}, []int{8, 2, 7, 3, 5, 4}},
		{[]int{0, 1, 3, 4, 6, 7, 8, 0}, []int{0, 1, 3, 4, 6, 7, 8}},
	}
	for i := range tests {
		r := UniqueWithoutSort(tests[i].arg)
		if !reflect.DeepEqual(r, tests[i].result) {
			t.Fatal(i, r, tests[i].result)
		}
	}
}
