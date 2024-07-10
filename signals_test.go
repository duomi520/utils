package utils

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestSignalState(t *testing.T) {
	p := NewProcessor()
	s := p.SignalState(1314)
	time.Sleep(time.Millisecond)
	if p.GetSignalState(s).(int) != 1314 {
		t.Error(p.GetSignalState(s))
	}
	p.Close()
	a := p.SignalState(2324)
	time.Sleep(time.Millisecond)
	if p.GetSignalState(a) != nil {
		t.Error(p.GetSignalState(a))
	}
}
func TestEffector(t *testing.T) {
	p := NewProcessor()
	defer p.Close()
	s1 := p.SignalState(1)
	s2 := p.SignalState(2)
	f := func(a *Processor) {
		fmt.Printf("The count is %v and %v\n", GetSignal(a, s1), GetSignal(a, s2))
	}
	e := p.Effector(f, s1, s2)
	p.SetSignalState(s1, 1314)
	time.Sleep(time.Millisecond)
	p.SetSignalState(s2, 520)
	time.Sleep(time.Millisecond)
	p.RemoveEffector(e)
	p.SetSignalState(s1, 100)
	time.Sleep(time.Millisecond)
}

// The count is 1314 and 2
// The count is 1314 and 520
func TestComputed(t *testing.T) {
	p := NewProcessor()
	defer p.Close()
	firstName := p.SignalState("John")
	lastName := p.SignalState("Smith")
	f1 := func(a *Processor) any {
		return GetSignal(a, firstName).(string) + "." + GetSignal(a, lastName).(string)
	}
	c1 := p.Computed(f1, firstName, lastName)
	if !strings.EqualFold(p.GetComputerState(c1).(string), "John.Smith") {
		t.Fatal(p.GetComputerState(c1))
	}
	f2 := func(a *Processor) any {
		return "You name is " + GetComputer(a, c1).(string)
	}
	c2 := p.Computed(f2, c1)
	p.SetSignalState(firstName, "Joke")
	if !strings.EqualFold(p.GetComputerState(c1).(string), "Joke.Smith") {
		t.Fatal(p.GetComputerState(c1))
	}
	if !strings.EqualFold(p.GetComputerState(c2).(string), "You name is Joke.Smith") {
		t.Fatal(p.GetComputerState(c2))
	}
	err := p.RemoveComputer(c1)
	if err == nil {
		t.Fatal("错误")
	}
	err = p.RemoveComputer(c2)
	if err != nil {
		t.Fatal(err.Error())
	}
	if p.GetComputerState(c2) != nil {
		t.Fatal(p.GetComputerState(c2))
	}
	p.SetSignalState(firstName, "Mike")
	if !strings.EqualFold(p.GetComputerState(c1).(string), "Mike.Smith") {
		t.Fatal(p.GetComputerState(c1))
	}
}
