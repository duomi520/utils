package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestSignal(t *testing.T) {
	u := NewUniverse()
	s := NewSignal(1314, nil)
	u.Run()
	if s.Load().(int) != 1314 {
		t.Error(s.Load())
	}
	u.Close()
	s.Set(u, 2321)
	time.Sleep(time.Millisecond)
	if s.Load().(int) != 1314 {
		t.Error(s.Load())
	}
	fmt.Println(s.Load())
}

// 1314

func TestComputed(t *testing.T) {
	u := NewUniverse()
	defer u.Close()
	firstName := NewSignal("John", nil)
	lastName := NewSignal("Smith", nil)
	f1 := func(a ...Parent) any {
		return a[0].Load().(string) + "." + a[1].Load().(string)
	}
	c1 := NewComputer(nil, f1, firstName, lastName)
	f2 := func(a ...Parent) any {
		return "You name is " + a[0].Load().(string)
	}
	c2 := NewComputer(nil, f2, c1)
	u.Run()
	c2.Operate(u)
	time.Sleep(time.Millisecond)
	fmt.Println(c2.Load())
	firstName.Set(u, "Joke")
	time.Sleep(time.Millisecond)
	fmt.Println(c2.Load())
	firstName.Set(u, "Mike")
	time.Sleep(time.Millisecond)
	fmt.Println(c1.Load())
}

//You name is John.Smith
//You name is Joke.Smith
//Mike.Smith

func TestEffector(t *testing.T) {
	u := NewUniverse()
	defer u.Close()
	f := func(a any) {
		fmt.Printf("The number is %v \n", a)
	}
	s1 := NewSignal(1314, f)
	s2 := NewSignal(520, f)
	fc := func(a ...Parent) any {
		return a[0].Load().(int)*1000 + a[1].Load().(int)
	}
	c := NewComputer(f, fc, s1, s2)
	u.Run()
	c.Operate(u)
	time.Sleep(time.Millisecond)
	s1.Set(u, 100)
	c.Operate(u)
	time.Sleep(time.Millisecond)
}

//The number is 1314520
//The number is 100
//The number is 100520
