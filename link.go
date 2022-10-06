package utils

import (
	_ "unsafe"
)

//go:linkname FastRand runtime.fastrand
//后续go版本无法预测不变
func FastRand() uint32

//go:linkname Nanotime runtime.nanotime1
//后续go版本无法预测不变
func Nanotime() int64

// https://mp.weixin.qq.com/s/IG4HRjU-pOeaKBZ1ZRSiSQ
