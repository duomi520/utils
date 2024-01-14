package utils

import (
	"sync/atomic"
	"time"
)

// 定义魔改版雪花算法的参数
// highest 1 bit: always 0
// next   10 bit: workerId
// next   41 bit: timestamp
// lowest 12 bit: sequence
const (
	//最大1023
	MaxWorkNumber      = 1023
	WorkLeftShift      = uint(12 + 41)
	TimestampLeftShift = uint(12)
)

// SnowFlakeIDPlus 魔改版雪花算法,缺失原算法包含的时间戳信息
type SnowFlakeIDPlus struct {
	//时间戳启动计算时间零点
	systemCenterStartupTime int64
	//10bit的工作机器id
	workID        int64
	lastTimestamp int64
}

// NewSnowFlakeIDPlus 新建
func NewSnowFlakeIDPlus(id int64, startupTime int64) *SnowFlakeIDPlus {
	if id < 0 || id >= MaxWorkNumber || startupTime < 0 {
		return nil
	}
	s := &SnowFlakeIDPlus{
		systemCenterStartupTime: startupTime,
		workID:                  id,
		lastTimestamp:           (id << WorkLeftShift) | ((time.Now().UnixNano()-startupTime)/int64(time.Millisecond))<<TimestampLeftShift,
	}
	return s
}

// NextID 取得 id.
func (s *SnowFlakeIDPlus) NextID() int64 {
	return atomic.AddInt64(&s.lastTimestamp, 1)
}

// GetWorkID 根据id计算工作机器id
func GetWorkID(id int64) int64 {
	temp := id >> WorkLeftShift
	//1111111111   10bit
	return temp & 1023
}

// GetWorkID 取得工作机器id
func (s *SnowFlakeIDPlus) GetWorkID() int64 {
	return s.workID
}

// 改良版的雪花算法  https://zhuanlan.zhihu.com/p/648460337
