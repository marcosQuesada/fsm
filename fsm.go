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
	mutex        *sync.RWMutex
}

func NewFSM(state State, states []State) *FSM {
	s := make(map[StateType]State, 0)
	for _, ss := range states {
		s[ss.GetType()] = ss
	}

	return &FSM{
		currentState: state,
		states:       s,
		mutex:        &sync.RWMutex{},
	}
}

func (f *FSM) GetCurrentState() State {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	return f.currentState
}

func (f *FSM) Boot() {
	f.mutex.Lock()
	s := f.currentState
	f.mutex.Unlock()

	s.Enable()

	t := f.GetTransitionReady()
	if t == nil {
		return
	}

	f.DoTransition(t)
}

func (f *FSM) GetTransitionReady() *Transition {
	f.mutex.Lock()
	s := f.currentState
	f.mutex.Unlock()

	for _, t := range s.Transitions() {

		if f.PassGuards(t) {
			return t
		}
	}

	return nil
}

func (f *FSM) DoTransition(t *Transition) {
	f.mutex.Lock()
	s := f.currentState
	f.mutex.Unlock()

	s.Disable()

	s, ok := f.states[t.To]
	if !ok {
		log.Print("Unexpected error, State not found ", t.To)
		return
	}

	f.mutex.Lock()
	f.currentState = s
	f.mutex.Unlock()

	s.Enable()
}

func (f *FSM) Terminate() {
	f.mutex.Lock()
	s := f.currentState
	f.mutex.Unlock()

	s.Disable()
}

func (f *FSM) PassGuards(t *Transition) bool {
	for _, g := range t.Guards {
		if !g() {
			return false
		}
	}

	return true
}
