package utils

import (
	"errors"
	"sync"
	"testing"
	"time"
)

//
// 工作组中心分配工作机器id
//

//ErrFailureAllocID 定义错误
var ErrFailureAllocID = errors.New("utils.IdWorker.GetId|工作机器id耗尽。")

//IDWorker 工作组用于分配工作机器id
type IDWorker struct {
	//2018-6-1 00:00:00 UTC  ，时间戳启动计算时间零点
	SystemCenterStartupTime int64 
	queue                   []int
	queueMap                []int
	//消费
	consumption             int 
	//生产
	production              int 
	mutex                   sync.Mutex
}

//NewIDWorker 初始化
func NewIDWorker(st int64) *IDWorker {
	w := &IDWorker{
		SystemCenterStartupTime: st,
		queue:                   make([]int, MaxWorkNumber),
		queueMap:                make([]int, MaxWorkNumber),
		production:              0,
		consumption:             0,
	}
	for i := 0; i < MaxWorkNumber; i++ {
		w.queue[i] = i
		w.queueMap[i] = i
	}
	return w
}

//GetID 取id
func (w *IDWorker) GetID() (int, error) {
	n := 0
	var err error
	w.mutex.Lock()
	if w.queue[w.consumption] == -1 {
		n = -1
		err = ErrFailureAllocID
	} else {
		n = w.queue[w.consumption]
		w.queue[w.consumption] = -1
		w.consumption++
		if w.consumption >= MaxWorkNumber {
			w.consumption = 0
		}
		w.queueMap[n] = -1
		err = nil
	}
	w.mutex.Unlock()
	return n, err
}

//
//当工作机器与工作组中心的时间不同步时，释放后再利用的workID与之前释放的workID的snowflake id会重复，逻辑上产生bug。
//

//PutID 还id
func (w *IDWorker) PutID(n int) {
	if n < 0 || n >= MaxWorkNumber {
		return
	}
	w.mutex.Lock()
	if w.queueMap[n] != -1 || w.queue[w.production] != -1 {
		w.mutex.Unlock()
		return
	}
	w.queue[w.production] = n
	w.queueMap[n] = w.production
	w.production++
	if w.production >= MaxWorkNumber {
		w.production = 0
	}
	w.mutex.Unlock()
}

func TestIDWorker(t *testing.T) {
	// Test table
	idWorkerTests := []int{2, 3, 8, 10, -1, MaxWorkNumber + 1, -10} 
	// Verify table
	idWorkerVerify := []int{2, 3, 8, 10}                            
	w := NewIDWorker(time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano())
	for i := 0; i < MaxWorkNumber; i++ {
		f, err := w.GetID()
		if f != i || err != nil {
			t.Error(f, err, w)
		}
	}
	g0, err := w.GetID()
	if err == nil {
		t.Error(g0, err, w)
	}
	for _, l := range idWorkerTests {
		w.PutID(l)
	}
	for _, l := range idWorkerVerify {
		g, err := w.GetID()
		if l != g || err != nil {
			t.Error(g, err, w)
		}

	}
}
func TestNextID(t *testing.T) {
	var v [30000]int64
	s1 := NewSnowFlakeID(1, time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano())
	for i := 0; i < 10000; i++ {
		v[i], _ = s1.NextID()
	}
	time.Sleep(time.Millisecond)
	s2 := NewSnowFlakeID(1, time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano())
	for i := 10000; i < 20000; i++ {
		v[i], _ = s2.NextID()
	}
	s3 := NewSnowFlakeID(3, time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano())
	for i := 20000; i < 30000; i++ {
		v[i], _ = s3.NextID()
	}
	//验证
	for i := 0; i < (30000 - 1); i++ {
		if v[i] >= v[i+1] {
			t.Error("失败：i:", i, "v[i]:", v[i], "v[i+1]:", v[i+1])
			t.FailNow()
		}
	}
}

func TestGetWorkID(t *testing.T) {
	s1 := NewSnowFlakeID(1, time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano())
	s2 := NewSnowFlakeID(555, time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano())
	s3 := NewSnowFlakeID(1022, time.Date(2017, time.January, 1, 0, 0, 0, 0, time.UTC).UnixNano())
	n1, _ := s1.NextID()
	n2, _ := s2.NextID()
	n3, _ := s3.NextID()
	if GetWorkID(n1) != 1 {
		t.Error("失败:", n1, GetWorkID(n1))
	}
	if GetWorkID(n2) != 555 {
		t.Error("失败:", n2, GetWorkID(n2))
	}
	if GetWorkID(n3) != 1022 {
		t.Error("失败:", n3, GetWorkID(n3))
	}
}
