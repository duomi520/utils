package utils

import (
	_ "unsafe"
)

//go:linkname FastRand runtime.fastrand
func FastRand() uint32

//go:linkname Nanotime runtime.nanotime1
func Nanotime() int64

// https://mp.weixin.qq.com/s/IG4HRjU-pOeaKBZ1ZRSiSQ
