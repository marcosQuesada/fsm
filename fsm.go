package fsm

import (
	"sync"
	"log"
)

type StateType string

type State interface {
	GetType() StateType
	Enable() error
	Disable() error
	Transitions() []*Transition
}

type Transition struct {
	To     StateType
	Guards []Guard
}

type Guard func() bool

type FSM struct {
	currentState State
	states       map[StateType]State
	mutex        *sync.Mutex
}

func NewFSM(state State, states []State) *FSM {
	s := make(map[StateType]State, 0)
	for _, ss := range states {
		s[ss.GetType()] = ss
	}
	return &FSM{
		currentState: state,
		states:       s,
		mutex:        &sync.Mutex{},
	}
}

func (f *FSM) GetCurrentState() State {
	return f.currentState
}

func (f *FSM) Boot() {
	f.currentState.Enable()

	t := f.GetTransitionReady()
	if t == nil {
		return
	}

	f.DoTransition(t)
}

func (f *FSM) GetTransitionReady() *Transition {
	for _, t := range f.currentState.Transitions() {

		if f.PassGuards(t) {
			return t
		}
	}

	return nil
}

func (f *FSM) DoTransition(t *Transition) {
	f.currentState.Disable()

	s, ok := f.states[t.To]
	if !ok {
		log.Print("Unexpected error, State not found ", t.To)
		return
	}

	f.mutex.Lock()
	f.currentState = s
	f.mutex.Unlock()

	f.currentState.Enable()
}

func (f *FSM) Terminate() {
	f.currentState.Disable()
}

func (f *FSM) PassGuards(t *Transition) bool {
	for _, g := range t.Guards {
		if !g() {
			return false
		}
	}

	return true
}
