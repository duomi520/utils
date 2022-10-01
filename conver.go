package utils

import (
	"unsafe"
)

//StringToBytes String转Bytes 转后不要修改Bytes
func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

//BytesToString Bytes转String 转后不要修改Bytes
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

//IntegerEqual 判断两个整数切片的内容是否完全相同
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

//DIGITS 数字串
var DIGITS = [10]byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

//UItoa 非负整数转String
func UItoa[T uint | uint8 | uint16 | uint32 | uint64](d T) string {
	if d == 0 {
		return "0"
	}
	var b [20]byte
	n := 0
	for j := 19; j >= 0 && d > 0; j-- {
		b[j] = DIGITS[d%10]
		d /= 10
		n++
	}
	return string(b[20-n:])
}

type Integer16 interface {
	int16 | uint16
}

//BytesToInteger16 切片转 int16 | uint16 little_endian
func BytesToInteger16[T Integer16](b []byte) T {
	_ = b[1]
	return T(b[0]) | T(b[1])<<8
}

//CopyInteger16 将 int16 | uint16 加入切片
func CopyInteger16[T Integer16](dst []byte, n T) {
	var c [2]byte
	c[0] = byte(n)
	c[1] = byte(n >> 8)
	copy(dst[:], c[:])
}

type Integer32 interface {
	int32 | uint32
}

//BytesToInteger32 切片转 int32 | uint32  little_endian
func BytesToInteger32[T Integer32](b []byte) T {
	_ = b[3]
	return T(b[0]) | T(b[1])<<8 | T(b[2])<<16 | T(b[3])<<24
}

//CopyInteger32 将 int32 | uint32 加入切片
func CopyInteger32[T Integer32](dst []byte, n T) {
	var c [4]byte
	c[0] = byte(n)
	c[1] = byte(n >> 8)
	c[2] = byte(n >> 16)
	c[3] = byte(n >> 24)
	copy(dst[:], c[:])
}

type Integer64 interface {
	int64 | uint64
}

//CopyInteger64 将 int64 | uint64 加入切片 little_endian
func CopyInteger64[T Integer64](dst []byte, n T) {
	var c [8]byte
	c[0] = byte(n)
	c[1] = byte(n >> 8)
	c[2] = byte(n >> 16)
	c[3] = byte(n >> 24)
	c[4] = byte(n >> 32)
	c[5] = byte(n >> 40)
	c[6] = byte(n >> 48)
	c[7] = byte(n >> 56)
	copy(dst[:], c[:])
}

//BytesToInteger64 切片转 int64 | uint64
func BytesToInteger64[T Integer64](b []byte) T {
	_ = b[7]
	return T(b[0]) | T(b[1])<<8 | T(b[2])<<16 | T(b[3])<<24 | T(b[4])<<32 | T(b[5])<<40 | T(b[6])<<48 | T(b[7])<<56
}
