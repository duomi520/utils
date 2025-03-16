package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestUniverse(t *testing.T) {
	u := NewUniverse()
	defer u.Close()
	s := u.NewSignal(1314, func(a any) {
		fmt.Println(a)
	})
	u.Run()
	fmt.Println(u.signalSet[-s].value)
	u.SetSignal(s, 2321)
	time.Sleep(time.Millisecond)
}

// 1314
// 2321

func TestComputed(t *testing.T) {
	show := func(a any) {
		fmt.Println(a)
	}
	u := NewUniverse()
	defer u.Close()
	firstName := u.NewSignal("John", nil)
	lastName := u.NewSignal("Smith", nil)
	f1 := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetSignal(firstName).(string) + "." + u.GetSignal(lastName).(string)
		}, []int{firstName, lastName}
	}
	c1 := u.NewComputer(nil, f1)
	f2 := func(u *Universe) (func() any, []int) {
		return func() any {
			return "You name is " + u.GetComputer(c1).(string)
		}, []int{c1}
	}
	c2 := u.NewComputer(nil, f2)
	u.Run()
	u.Operate(c2, show)
	time.Sleep(time.Millisecond)
	u.SetSignal(firstName, "Joke")
	u.Operate(c2, show)
	time.Sleep(time.Millisecond)
	u.SetSignal(firstName, "Mike")
	time.Sleep(time.Millisecond)
	u.Operate(c1, show)
	time.Sleep(time.Millisecond)
}

//You name is John.Smith
//You name is Joke.Smith
//Mike.Smith

func TestEffector(t *testing.T) {
	u := NewUniverse()
	defer u.Close()
	watch := func(a any) {
		fmt.Printf("The number is %v \n", a)
	}
	s1 := u.NewSignal(1314, watch)
	s2 := u.NewSignal(520, watch)
	fc := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetSignal(s1).(int)*1000 + u.GetSignal(s2).(int)
		}, []int{s1, s2}
	}
	c := u.NewComputer(watch, fc)
	u.Run()
	u.Operate(c, nil)
	time.Sleep(time.Millisecond)
	u.SetSignal(s1, 100)
	u.Operate(c, nil)
	time.Sleep(time.Millisecond)
}

//The number is 1314520
//The number is 100
//The number is 100520
