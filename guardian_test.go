package utils

import (
	"log"
	"testing"
	"time"
)

func testJob1() bool {
	log.Println("Job1 do")
	return true
}
func testJob2() bool {
	log.Println("Job2 do")
	return false
}

var testJob3Time int64

func testJob3() bool {
	log.Println(time.Duration((Nanotime() - testJob3Time)))
	testJob3Time = Nanotime()
	return false
}
func TestAddJob(t *testing.T) {
	logger, _ := NewWLogger(DebugLevel, "")
	g := NewGuardian(80*time.Millisecond, logger)
	g.AddJob(testJob1)
	g.AddJob(testJob2)
	g.AddJob(testJob2)
	time.Sleep(250 * time.Millisecond)
	g.Release()
	time.Sleep(250 * time.Millisecond)
}

/*
[Debug] 2022-11-26 14:19:20 NewGuardian: 启用定时器
2022/11/26 14:19:20 Job1 do
2022/11/26 14:19:20 Job2 do
2022/11/26 14:19:20 Job2 do
2022/11/26 14:19:20 Job2 do
2022/11/26 14:19:20 Job2 do
2022/11/26 14:19:20 Job2 do
[Debug] 2022-11-26 14:19:20 NewGuardian: 定时器关闭
*/
func TestGuardian(t *testing.T) {
	logger, _ := NewWLogger(DebugLevel, "")
	g := NewGuardian(100*time.Millisecond, logger)
	testJob3Time = Nanotime()
	g.AddJob(testJob3)
	time.Sleep(1000 * time.Millisecond)
	g.Release()
	time.Sleep(250 * time.Millisecond)
}

/*
[Debug] 2022-12-07 20:54:05 NewGuardian: 启用定时器
2022/12/07 20:54:05 9.232ms
2022/12/07 20:54:05 112.4186ms
2022/12/07 20:54:05 108.7808ms
2022/12/07 20:54:05 107.9704ms
2022/12/07 20:54:05 109.9366ms
2022/12/07 20:54:05 110.2107ms
2022/12/07 20:54:05 106.5798ms
2022/12/07 20:54:05 110.6951ms
2022/12/07 20:54:05 108.169ms
2022/12/07 20:54:06 108.1498ms
[Debug] 2022-12-07 20:54:06 NewGuardian: 定时器关闭
*/
