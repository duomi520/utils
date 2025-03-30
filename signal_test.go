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

func TestComputer(t *testing.T) {
	show := func(a any) {
		fmt.Println(a)
	}
	u := NewUniverse()
	defer u.Close()
	firstName := u.NewSignal("John", nil)
	lastName := u.NewSignal("Smith", nil)
	f1 := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetString(firstName) + "." + u.GetString(lastName)
		}, []int{firstName, lastName}
	}
	c1 := u.NewComputer(nil, f1)
	f2 := func(u *Universe) (func() any, []int) {
		return func() any {
			return "You name is " + u.GetString(c1)
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
			return u.GetInt(s1)*1000 + u.GetInt(s2)
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

func TestCompute(t *testing.T) {
	show := func(u *Universe, key string, a any) {
		for v := range u.signalSet[1:] {
			fmt.Printf(" %v ", u.signalSet[1:][v].value)
		}
		fmt.Printf("|")
		for v := range u.computerSet[1:] {
			fmt.Printf(" %v ", u.computerSet[1:][v].value)
		}
		fmt.Printf("|")
		fmt.Printf(" %s = %v \n", key, a)
	}
	u := NewUniverse()
	defer u.Close()
	s1 := u.NewSignal(1, func(a any) {
		fmt.Printf(" s1 = %v \n", a)
	})
	s2 := u.NewSignal(2, func(a any) {
		fmt.Printf(" s2 = %v \n", a)
	})
	s3 := u.NewSignal(3, func(a any) {
		fmt.Printf(" s3 = %v \n", a)
	})
	f11 := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetInt(s1) + u.GetInt(s2)
		}, []int{s1, s2}
	}
	f12 := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetInt(s2) * 2
		}, []int{s2}
	}
	c11 := u.NewComputer(func(a any) { show(u, "c11", a) }, f11)
	c12 := u.NewComputer(func(a any) { show(u, "c12", a) }, f12)
	f21 := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetInt(c11) * 2
		}, []int{c11}
	}
	f22 := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetInt(c11) + u.GetInt(c12)
		}, []int{c11, c12}
	}
	c21 := u.NewComputer(func(a any) { show(u, "c21", a) }, f21)
	c22 := u.NewComputer(func(a any) { show(u, "c22", a) }, f22)
	f31 := func(u *Universe) (func() any, []int) {
		return func() any {
			return u.GetInt(c22) + u.GetInt(s3)
		}, []int{c22, s3}
	}
	c31 := u.NewComputer(func(a any) {show(u, "c31", a) }, f31)
	u.Run()
	u.Operate(c11, nil)
	u.Operate(c21, nil)
	time.Sleep(time.Millisecond)
	u.SetSignal(s1, 6)
	u.Operate(c31, nil)
	time.Sleep(time.Millisecond)
	u.SetSignal(s1, 4)
	u.SetSignal(s3, 8)
	time.Sleep(time.Millisecond)
	u.Operate(c31, nil)
	time.Sleep(time.Millisecond)
	u.Operate(c21, nil)
	time.Sleep(time.Millisecond)
}

/*
 1  2  3 | 3  <nil>  <nil>  <nil>  <nil> | c11 = 3
 1  2  3 | 3  <nil>  6  <nil>  <nil> | c21 = 6
 s1 = 6
 6  2  3 | 8  <nil>  6  <nil>  <nil> | c11 = 8
 6  2  3 | 8  4  6  <nil>  <nil> | c12 = 4
 6  2  3 | 8  4  6  12  <nil> | c22 = 12
 6  2  3 | 8  4  6  12  15 | c31 = 15
 s1 = 4
 s3 = 8
 4  2  8 | 6  4  6  12  15 | c11 = 6
 4  2  8 | 6  4  6  10  15 | c22 = 10
 4  2  8 | 6  4  6  10  18 | c31 = 18
 4  2  8 | 6  4  12  10  18 | c21 = 12
*/
