package fsm

type FiniteStateMachine struct {
	name      string
	curState  *State
	stateList map[int]*State
}

func NewFiniteStateMachine(name string) *FiniteStateMachine {
	return &FiniteStateMachine{
		name:      name,
		curState:  nil,
		stateList: make(map[int]*State)}
}

func (f *FiniteStateMachine) GetState(id int) *State {
	return f.stateList[id]
}

func (f *FiniteStateMachine) AddState(id int) *State {
	s := f.GetState(id)
	if s == nil {
		s = NewState(id)
		f.stateList[id] = s
	}
	return s
}

func (f *FiniteStateMachine) AddTransition(id int, in int, out int) bool {
	s := f.AddState(id)
	if s == nil {
		return false
	}
	return s.addTransition(in, out)
}

func (f *FiniteStateMachine) DeleteTransition(id int, in int) {
	s := f.GetState(id)
	if s == nil {
		return
	}
	s.deleteTransition(in)
	if s.getStateCount() == 0 {
		delete(f.stateList, id)
	}
}

func (f *FiniteStateMachine) CanTransitState(id int, in int) bool {
	s := f.GetState(id)
	if s == nil {
		return false
	}
	return s.findOutputState(in)
}

func (f *FiniteStateMachine) SetCurrentState(id int) bool {
	s := f.GetState(id)
	if s == nil {
		return false
	}

	f.curState = s
	return true
}

func (f *FiniteStateMachine) TransitState(in int) bool {
	if f.curState == nil {
		return false
	}

	var out int
	if !f.curState.getOutState(in, &out) {
		return false
	}
	return f.SetCurrentState(out)
}

func (f *FiniteStateMachine) getCurrentStateId() int {
	if f.curState == nil {
		return InvalidStateID
	}
	return f.curState.ID
}

func (f *FiniteStateMachine) FindOutputStateId(in int) int {
	s := f.GetState(f.getCurrentStateId())
	if s == nil {
		return InvalidStateID
	}
	var out int
	if !s.getOutState(in, &out) {
		return InvalidStateID
	}
	return out
}
