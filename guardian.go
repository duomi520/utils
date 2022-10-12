package utils

import (
	"context"
	"errors"
	"sync"
	"time"
)

//Guardian 定时任务协程
type Guardian struct {
	ctx    context.Context
	Period time.Duration
	//返回true 退出定时任务
	addChan   chan func() bool
	closeOnce sync.Once
	stopChan  chan struct{}
}

//NewGuardian 新建
func NewGuardian(ctx context.Context) *Guardian {
	var g = Guardian{
		ctx:      ctx,
		Period:   5 * time.Second,
		addChan:  make(chan func() bool, 16),
		stopChan: make(chan struct{}),
	}
	go g.Run()
	return &g
}

//Release 释放
func (g *Guardian) Release() {
	g.closeOnce.Do(func() {
		close(g.stopChan)
	})
}

//AddJob 加入
func (g *Guardian) AddJob(f func() bool) error {
	select {
	case g.addChan <- f:
	default:
		return errors.New("Guardian 加入任务失败")
	}
	return nil
}

//Run 协程
func (g *Guardian) Run() {
	ticker := time.NewTicker(g.Period)
	var chain []func() bool
	defer func() {
		ticker.Stop()
		g.Release()
	}()
	for {
		select {
		case <-ticker.C:
			i := 0
			lenght := len(chain)
			for i < lenght {
				if chain[i]() {
					if i < (lenght - 1) {
						copy(chain[i:], chain[i+1:])
					}
					lenght--
				}
				i++
			}
			chain = chain[:lenght]
		case f := <-g.addChan:
			chain = append(chain, f)
		case <-g.ctx.Done():
			return
		case <-g.stopChan:
			return
		}
	}

}
