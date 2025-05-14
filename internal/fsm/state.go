package fsm

const (
	InvalidStateID int = -1
)

type State struct {
	ID          int
	transitions map[int]int
}

func NewState(sid int) *State {
	return &State{sid, make(map[int]int)}
}

func (s *State) addTransition(in int, out int) bool {
	if s.findOutputState(in) {
		return false
	}

	s.transitions[in] = out
	return true
}

func (s *State) deleteTransition(in int) {
	delete(s.transitions, in)
}

func (s *State) getOutState(in int) (int, bool) {
	v, ok := s.transitions[in]
	if !ok {
		return InvalidStateID, false
	}
	return v, true
}

func (s *State) findOutputState(in int) bool {
	if _, ok := s.transitions[in]; ok {
		return true
	}
	return false
}

func (s *State) getStateCount() int {
	return len(s.transitions)
}
