package utils


import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)


func TestBroadcast(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	o := NewOptions()
	s := NewService(o)
	err := s.TCPServer(context.TODO(), ":4567")
	if err != nil {
		t.Fatal(err.Error())
	}
	topic := s.RegisterTopic("room")
	c1, err := NewTCPClient(context.TODO(), "127.0.0.1:4567", o)
	if err != nil {
		t.Fatal(err.Error())
	}
	c2, err := NewTCPClient(context.TODO(), "127.0.0.1:4567", o)
	if err != nil {
		t.Fatal(err.Error())
	}
	f := func(data []byte) error {
		fmt.Println(string(data))
		wg.Done()
		return nil
	}
	err = c1.Subscribe("room", f)
	defer c1.Unsubscribe("room")
	if err != nil {
		t.Fatal(err.Error())
	}
	err = c2.Subscribe("room", f)
	defer c2.Unsubscribe("room")
	if err != nil {
		t.Fatal(err.Error())
	}
	time.Sleep(400 * time.Millisecond)
	err = topic.Broadcast([]byte("Good!"))
	if err != nil {
		t.Fatal(err.Error())
	}
	wg.Wait()
}

