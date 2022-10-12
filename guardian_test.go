package utils

import (
	"context"
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
	ctx, ctxExitFunc := context.WithCancel(context.Background())
	g := NewGuardian(ctx)
	g.Period = 80 * time.Millisecond
	go g.Run()
	g.AddJob(testJob1)
	g.AddJob(testJob2)
	g.AddJob(testJob2)
	time.Sleep(250 * time.Millisecond)
	ctxExitFunc()
	time.Sleep(250 * time.Millisecond)
}
