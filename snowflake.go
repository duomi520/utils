package utils

import (
	"errors"
	"sync"
	"time"
)

// 定义snowflake的参数
const (
	//最大1023
	MaxWorkNumber = 1023
	//理论最大4095
	MaxSequenceNumber  = 4000
	WorkLeftShift      = uint(12)
	TimestampLeftShift = uint(22)
)

// ErrMachineTimeUnSynchronize 定义错误
var ErrMachineTimeUnSynchronize = errors.New("机器时钟后跳，timeGen()生成时间戳，早于SnowFlakeId记录的时间戳")

//
//各个工作站生成Twitter-Snowflake算法的ID
//

// SnowFlakeID 工作站
type SnowFlakeID struct {
	//时间戳启动计算时间零点
	systemCenterStartupTime int64
	//41bit的时间戳，仅支持69.7年
	lastTimestamp int64
	//10bit的工作机器id
	workID int64
	//12bit的序列号
	sequence int64
	mutex    sync.Mutex
}

// NewSnowFlakeID 工作组
func NewSnowFlakeID(id int64, startupTime int64) *SnowFlakeID {
	if id < 0 || id >= MaxWorkNumber || startupTime < 0 {
		return nil
	}
	s := &SnowFlakeID{
		systemCenterStartupTime: startupTime / int64(time.Millisecond),
		lastTimestamp:           timeGen(),
		workID:                  id << WorkLeftShift,
		sequence:                0,
	}
	return s
}

// NextID 取得 snowflake id.
func (s *SnowFlakeID) NextID() (int64, error) {
	timestamp := timeGen()
	s.mutex.Lock()
	if timestamp < s.lastTimestamp {
		s.mutex.Unlock()
		return 0, ErrMachineTimeUnSynchronize
	}
	if timestamp == s.lastTimestamp {
		s.sequence = s.sequence + 1
		if s.sequence > MaxSequenceNumber {
			//效率不高，貌似影响不大
			time.Sleep(time.Millisecond)
			timestamp++
			s.sequence = 0
		}
	} else {
		s.sequence = 0
	}
	s.lastTimestamp = timestamp
	id := ((timestamp - s.systemCenterStartupTime) << TimestampLeftShift) | s.workID | s.sequence
	s.mutex.Unlock()
	return id, nil
}

// timeGen 取得time.Now() unix 毫秒.
func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// GetWorkID 取得工作机器id
func GetWorkID(id int64) int64 {
	temp := id >> WorkLeftShift
	//1111111111   10bit
	return temp & 1023
}

// GetWorkID 取得工作机器id
func (s *SnowFlakeID) GetWorkID() int64 {
	return s.workID >> WorkLeftShift
}
