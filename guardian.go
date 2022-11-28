package utils

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

//Guardian 批量处理定时任务协程
type Guardian struct {
	Period    int64
	lock      sync.Mutex
	chain     []func() bool
	closeOnce sync.Once
	stopChan  chan struct{}
	logger    ILogger
}

//NewGuardian 新建批量处理定时器
func NewGuardian(period time.Duration, l ILogger) *Guardian {
	var g = Guardian{
		Period:   int64(period),
		stopChan: make(chan struct{}),
		logger:   l,
	}
	go func() {
		g.logger.Debug("NewGuardian: 启用定时器")
		for {
			select {
			case <-g.stopChan:
				g.logger.Debug("NewGuardian: 定时器关闭")
				return
			default:
				i := 0
				start := Nanotime()
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
					g.logger.Warn(fmt.Sprintf("NewGuardian：运行任务耗时 %d 超过设定周期 %d ", since, g.Period))
				}
			}
		}
	}()
	return &g
}

//Release 释放
func (g *Guardian) Release() {
	g.closeOnce.Do(func() {
		close(g.stopChan)
	})
}

//AddJob 加入定时任务，任务返回true 任务将退出定时运行，任务不可长时间阻塞
func (g *Guardian) AddJob(f func() bool) error {
	select {
	case <-g.stopChan:
		return errors.New("Guardian.AddJob：定时器已关闭")
	default:
		g.lock.Lock()
		g.chain = append(g.chain, f)
		g.lock.Unlock()
	}
	return nil
}
