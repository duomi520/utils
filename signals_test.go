package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestRemoveDuplicates(t *testing.T) {
	nums := []int{8, 2, 7, 3, 5, 3, 4, 5}
	uniqueNums := removeDuplicates(nums)
	fmt.Println(uniqueNums)
}

// [2 3 4 5 7 8]

func TestUniverse(t *testing.T) {
	u := NewUniverse()
	defer u.Close()
	s := u.NewSignal(1314, func(a any) {
		fmt.Println(a)
	})
	u.Run()
	fmt.Println(u.signalValue[0])
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
	f1 := func() (func() any, []Parent) {
		return func() any {
			return firstName.String() + "." + lastName.String()
		}, []Parent{firstName, lastName}
	}
	c1 := NewComputer(u, nil, f1)
	f2 := func() (func() any, []Parent) {
		return func() any {
			return "You name is " + c1.String()
		}, []Parent{c1}
	}
	c2 := NewComputer(u, nil, f2)
	u.Run()
	u.Operate(c2, show)
	time.Sleep(time.Millisecond)
	u.SetSignal(firstName, "Joke")
	u.Operate(c2, show)
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
	fc := func() (func() any, []Parent) {
		return func() any {
			return s1.Int()*1000 + s2.Int()
		}, []Parent{s1, s2}
	}
	c := NewComputer(u, watch, fc)
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
