package utils

import (
	"strconv"
	"testing"
	"time"
)

func TestWLogger(t *testing.T) {
	l, err := NewWLogger(DebugLevel, "")
	if err != nil {
		t.Fatal(err)
	}
	l.Debug("Hello World！")
	l.Debug("")
	l.Info("Logger ", "Debug")
	l.Warn("！@#￥%……&*（）")
	l.Error("错误")
	l.Close()
}
func TestLoggerFile(t *testing.T) {
	l, err := NewWLogger(DebugLevel, "log")
	if err != nil {
		t.Fatal(err)
	}
	l.SetLogFileSize(200, 5)
	b := "<----------------------------------------------------------------------------------------->"
	for i := 0; i < 5; i++ {
		time.Sleep(time.Second)
		l.Info(strconv.Itoa(i), b, strconv.Itoa(i))
	}
	l.Close()
}
