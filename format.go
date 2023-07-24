package utils

import "runtime"

func formatRecover() ([]byte, any) {
	if r := recover(); r != nil {
		const size = 65536
		buf := make([]byte, size)
		end := runtime.Stack(buf, false)
		if end > size {
			end = size
		}
		return buf[:end], r
	}
	return nil, nil
}
