package fsm

import (
	"testing"
	"log"
	"sync/atomic"
	"fmt"
)

func TestBasicFSMDefinition(t *testing.T) {
	var enables uint64 = 0
	var disables uint64 = 0
	var transitions uint64 = 0

	foo := &fooState{&enables, &disables, &transitions}
	bar := &barState{&enables, &disables, &transitions}
	fooBar := &fooBarState{&enables, &disables, &transitions}
	states := []State{foo, bar, fooBar}
	fsm := NewFSM(states[0], states)


	s := fsm.GetCurrentState()
	if s.GetType() != StateType("foo") {
		t.Error("Unexpected State")
	}

	fsm.Boot()

	tr := fsm.GetTransitionReady()
	if tr == nil {
		t.Error("Expected transition not found")
	}

	fsm.DoTransition(tr)

	if fsm.GetCurrentState().GetType() != StateType("fooBar") {
		t.Error("Unexpected State")
	}

	fsm.Terminate()
	e := atomic.LoadUint64(&enables)
	d := atomic.LoadUint64(&disables)
	trs := atomic.LoadUint64(&transitions)
	log.Println(fmt.Sprintf("Total enable %d disable %d transitions %d", e, d, trs))

	if e != 3 {
		t.Error("Unexpected Total Enable Operations")
	}

	if d != 3 {
		t.Error("Unexpected Total Disable Operations")
	}

	if trs != 2 {
		t.Error("Unexpected Total Transition Operations")
	}
}

type fooState struct {
	enables     *uint64
	disables    *uint64
	transitions *uint64
}

func (s *fooState) GetType() StateType{
	return StateType("foo")
}

func (s *fooState) Enable() error{
	log.Println("FooState Enable")
	atomic.AddUint64(s.enables, 1)

	return nil
}

func (s *fooState) Disable() error{
	log.Println("FooState Disable")
	atomic.AddUint64(s.disables, 1)
	return nil
}

func (s *fooState) Transitions() []*Transition{
	log.Println("FooState Transitions")
	atomic.AddUint64(s.transitions, 1)

	t := make([]*Transition, 0)
	t = append(t, &Transition{StateType("bar"), []Guard{s.done}})

	return t
}

func (s *fooState) done() bool {
	return true
}

type barState struct {
	enables     *uint64
	disables    *uint64
	transitions *uint64
}

func (s *barState) GetType() StateType{
	return StateType("bar")
}

func (s *barState) Enable() error{
	log.Println("BarState Enable")
	atomic.AddUint64(s.enables, 1)

	return nil
}

func (s *barState) Disable() error{
	log.Println("BarState Disable")
	atomic.AddUint64(s.disables, 1)

	return nil
}

func (s *barState) Transitions() []*Transition{
	log.Println("BarState Transitions")
	atomic.AddUint64(s.transitions, 1)

	t := make([]*Transition, 0)
	t = append(t, &Transition{StateType("fooBar"), []Guard{s.done}})

	return t
}

func (s *barState) done() bool {
	return true
}

type fooBarState struct {
	enables     *uint64
	disables    *uint64
	transitions *uint64
}

func (s *fooBarState) GetType() StateType{
	return StateType("fooBar")
}

func (s *fooBarState) Enable() error{
	log.Println("FooBarState Enable")
	atomic.AddUint64(s.enables, 1)

	return nil
}

func (s *fooBarState) Disable() error{
	log.Println("FooBarState Disable")
	atomic.AddUint64(s.disables, 1)

	return nil
}

func (s *fooBarState) Transitions() []*Transition{
	log.Println("fooBarState Transitions")
	atomic.AddUint64(s.transitions, 1)

	t := make([]*Transition, 0)
	t = append(t, &Transition{StateType("foo"), []Guard{s.done}})

	return t
}

func (s *fooBarState) done() bool {
	return true
}
