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
func TestGuardian(t *testing.T) {
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
