package utils

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

//Guardian 批量处理定时任务协程,精度不高
type Guardian struct {
	Period int64
	lock   sync.Mutex
	chain  []func() bool
	//退出标志 1-退出
	stopFlag int32
	logger   ILogger
}

//NewGuardian 新建批量处理定时器
//有误差，以ms计
func NewGuardian(period time.Duration, l ILogger) *Guardian {
	var g = Guardian{
		Period:   int64(period),
		stopFlag: 0,
		logger:   l,
	}
	go func() {
		g.logger.Debug("NewGuardian: 启用定时器")
		defer g.logger.Debug("NewGuardian: 定时器关闭")
		for {
			start := Nanotime()
			if atomic.LoadInt32(&g.stopFlag) == 1 {
				return
			}
			i := 0
			g.lock.Lock()
			lenght := len(g.chain)
			for i < lenght {
				if g.chain[i]() {
					if i < (lenght - 1) {
						copy(g.chain[i:], g.chain[i+1:])
					}
					lenght--
					g.chain = g.chain[:lenght]
				}
				i++
			}
			g.lock.Unlock()
			since := Nanotime() - start
			if since < g.Period {
				time.Sleep(time.Duration(g.Period - since))
			} else {
				g.logger.Warn(fmt.Sprintf("NewGuardian：运行任务耗时 %v 超过设定周期 %v ", time.Duration(since), time.Duration(g.Period)))
			}

		}
	}()
	return &g
}

//Release 释放
func (g *Guardian) Release() {
	atomic.StoreInt32(&g.stopFlag, 1)
}

//AddJob 加入定时任务，任务返回true 任务将退出定时运行，任务不可长时间阻塞
func (g *Guardian) AddJob(f func() bool) error {
	if atomic.LoadInt32(&g.stopFlag) == 1 {
		return nil
	}
	g.lock.Lock()
	g.chain = append(g.chain, f)
	g.lock.Unlock()
	return nil
}
