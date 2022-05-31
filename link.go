package utils

import (
	_ "unsafe"
)

//go:linkname FastRand runtime.fastrand
func FastRand() uint32

//go:linkname nanotime1 runtime.nanotime1
func nanotime1() int64

// https://mp.weixin.qq.com/s/IG4HRjU-pOeaKBZ1ZRSiSQ
