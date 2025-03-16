package utils

import (
	"fmt"
	"testing"
)

func TestRemoveDuplicates(t *testing.T) {
	nums := []int{8, 2, 7, 3, 5, 3, 4, 5}
	fmt.Println(RemoveDuplicates(nums))
}

// [2 3 4 5 7 8]

func TestUniqueWithoutSort(t *testing.T) {
	nums := []int{8, 2, 7, 3, 5, 3, 4, 5}
	fmt.Println(UniqueWithoutSort(nums))
}

// [8 2 7 3 5 4]
