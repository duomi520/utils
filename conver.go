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
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

type Integer16 interface {
	int16 | uint16
}

// BytesToInteger16 切片转 int16 | uint16 little_endian
func BytesToInteger16[T Integer16](b []byte) T {
	_ = b[1]
	return T(b[0]) | T(b[1])<<8
}

// CopyInteger16 将 int16 | uint16 加入切片
func CopyInteger16[T Integer16](dst []byte, n T) {
	_ = dst[1]
	dst[0] = byte(n)
	dst[1] = byte(n >> 8)
}

type Integer32 interface {
	int32 | uint32
}

// BytesToInteger32 切片转 int32 | uint32  little_endian
func BytesToInteger32[T Integer32](b []byte) T {
	_ = b[3]
	return T(b[0]) | T(b[1])<<8 | T(b[2])<<16 | T(b[3])<<24
}

// CopyInteger32 将 int32 | uint32 加入切片
func CopyInteger32[T Integer32](dst []byte, n T) {
	_ = dst[3]
	dst[0] = byte(n)
	dst[1] = byte(n >> 8)
	dst[2] = byte(n >> 16)
	dst[3] = byte(n >> 24)
}

type Integer64 interface {
	int64 | uint64
}

// CopyInteger64 将 int64 | uint64 加入切片 little_endian
func CopyInteger64[T Integer64](dst []byte, n T) {
	_ = dst[7]
	dst[0] = byte(n)
	dst[1] = byte(n >> 8)
	dst[2] = byte(n >> 16)
	dst[3] = byte(n >> 24)
	dst[4] = byte(n >> 32)
	dst[5] = byte(n >> 40)
	dst[6] = byte(n >> 48)
	dst[7] = byte(n >> 56)
}

// BytesToInteger64 切片转 int64 | uint64
func BytesToInteger64[T Integer64](b []byte) T {
	_ = b[7]
	return T(b[0]) | T(b[1])<<8 | T(b[2])<<16 | T(b[3])<<24 | T(b[4])<<32 | T(b[5])<<40 | T(b[6])<<48 | T(b[7])<<56
}

func offset(p unsafe.Pointer, n uintptr) unsafe.Pointer { return unsafe.Pointer(uintptr(p) + n) }
