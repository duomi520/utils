package utils

import (
	"cmp"
	"slices"
)

// RemoveDuplicates 排序去重
func RemoveDuplicates[T cmp.Ordered](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	// 先排序
	slices.Sort(s)
	// k用来记录不重复元素的索引位置
	k := 0
	for i := 1; i < len(s); i++ {
		if s[k] != s[i] {
			k++
			// 将不重复的元素移到数组前面
			s[k] = s[i]
		}
	}
	// 返回不重复元素的部分切片
	return s[:k+1]
}

// UniqueWithoutSort 去重但不排序
func UniqueWithoutSort[T comparable](s []T) []T {
	if len(s) == 0 {
		return nil
	}
	// 使用map来记录元素是否已经出现过
	seen := make(map[T]bool)
	var result []T
	for _, value := range s {
		// 如果这个值没有在map中出现过，就添加到结果切片中
		if _, ok := seen[value]; !ok {
			seen[value] = true
			result = append(result, value)
		}
	}
	return result
}
