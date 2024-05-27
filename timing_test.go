package utils

import (
	"container/heap"
	"fmt"
	"testing"
	"time"
)

func TestPriorityQueue(t *testing.T) {
	space := 100 * time.Millisecond
	now := time.Now()
	// 一些元素以及它们的优先级
	items := map[string]task{
		"3": {next: now.Add(3 * space)}, "2": {next: now.Add(2 * space)}, "4": {next: now.Add(4 * space)},
	}
	// 创建一个优先队列，并将上述元素放入到队列里面
	pq := make(priorityQueue, len(items))
	i := 0
	for _, priority := range items {
		pq[i] = priority
		pq[i].index = i
		i++
	}
	heap.Init(&pq)
	// 插入新元素，然后修改它的优先级
	heap.Push(&pq, task{next: now.Add(1 * space)})
	heap.Push(&pq, task{next: now.Add(10 * space)})
	// 以降序形式取出并打印队列中的所有元素
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(task)
		fmt.Printf("%v\n", item.next)
	}
}

/*
2024-05-26 17:26:45.5314494 +0800 CST m=+0.103569101
2024-05-26 17:26:45.6314494 +0800 CST m=+0.203569101
2024-05-26 17:26:45.7314494 +0800 CST m=+0.303569101
2024-05-26 17:26:45.8314494 +0800 CST m=+0.403569101
2024-05-26 17:26:46.4314494 +0800 CST m=+1.003569101
*/

func TestTimerRun(t *testing.T) {
	tr := NewTiming(nil)
	fa := func() time.Duration {
		fmt.Println("aa", time.Now())
		return 300 * time.Millisecond
	}
	fb := func() time.Duration {
		fmt.Println("bb", time.Now())
		return 0
	}
	fc := func() time.Duration {
		fmt.Println("cc", time.Now())
		return 500 * time.Millisecond
	}
	n := time.Now()
	tr.AddTask(n.Add(300*time.Millisecond), fa)
	tr.AddTask(n.Add(1000*time.Millisecond), fb)
	tr.AddTask(n.Add(500*time.Millisecond), fc)
	time.Sleep(3 * time.Second)
}

/*
aa 2024-05-27 20:02:16.6131111 +0800 CST m=+0.308639401
cc 2024-05-27 20:02:16.8129106 +0800 CST m=+0.508438901
aa 2024-05-27 20:02:16.9382378 +0800 CST m=+0.633766101
aa 2024-05-27 20:02:17.2480201 +0800 CST m=+0.943548401
bb 2024-05-27 20:02:17.3112233 +0800 CST m=+1.006751601
cc 2024-05-27 20:02:17.3267808 +0800 CST m=+1.022309101
aa 2024-05-27 20:02:17.558895 +0800 CST m=+1.254423301
cc 2024-05-27 20:02:17.8371953 +0800 CST m=+1.532723601
aa 2024-05-27 20:02:17.8679513 +0800 CST m=+1.563479601
aa 2024-05-27 20:02:18.1788282 +0800 CST m=+1.874356501
cc 2024-05-27 20:02:18.3487037 +0800 CST m=+2.044232001
aa 2024-05-27 20:02:18.4887739 +0800 CST m=+2.184302201
aa 2024-05-27 20:02:18.79801 +0800 CST m=+2.493538301
cc 2024-05-27 20:02:18.8595002 +0800 CST m=+2.555028501
aa 2024-05-27 20:02:19.1065649 +0800 CST m=+2.802093201
*/
