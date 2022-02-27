package utils

import (
	"runtime"
)

var defaultBallast []byte

//Ballast
func Ballast() {
	if len(defaultBallast) == 0 {
		//分配1G
		defaultBallast = make([]byte, 1*1024*1024*1024)
		//利用 runtime.KeepAlive 来保证 ballast 不会被 GC 给回收掉
		runtime.KeepAlive(defaultBallast)
	}
}

// https://www.cnblogs.com/457220157-FTD/p/15567442.html
