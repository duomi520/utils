package utils

import (
	"unsafe"
)

// StringToBytes String转Bytes 转后不要修改Bytes
func StringToBytes(s string) (b []byte) {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString Bytes转String 转后不要修改Bytes
func BytesToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

// IntegerEqual 判断两个整数切片的内容是否完全相同
func IntegerEqual[T Integer](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func offset(p unsafe.Pointer, n uintptr) unsafe.Pointer { return unsafe.Pointer(uintptr(p) + n) }
