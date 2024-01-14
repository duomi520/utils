package utils

import (
	_ "unsafe"
)

// 后续go版本无法预测不变
//
//go:linkname FastRand runtime.fastrand
func FastRand() uint32

// 后续go版本无法预测不变
//
//go:linkname Nanotime runtime.nanotime1
func Nanotime() int64

// https://mp.weixin.qq.com/s/IG4HRjU-pOeaKBZ1ZRSiSQ
