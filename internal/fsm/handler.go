package fsm

import (
	"sync"
	"time"
)

type Handler struct {
	fsm        *FiniteStateMachine
	curAction  ActionBase
	actionList []ActionBase
	Ch         chan int
}

func NewHandler(count int, name string) *Handler {
	return &Handler{
		fsm:        NewFiniteStateMachine(name),
		actionList: make([]ActionBase, count),
		Ch:         make(chan int),
	}
}

func (h *Handler) AddStateWithAction(id int, action ActionBase) bool {
	if action == nil {
		return false
	}
	h.actionList[id] = action
	h.fsm.AddState(id)
	return true
}

func (h *Handler) GetStateAction(id int) ActionBase {
	if !h.IsValidStateID(id) {
		return nil
	}
	return h.actionList[id]
}

func (h *Handler) AddTransition(id int, in int, out int, d float32) bool {
	a := h.GetStateAction(id)
	if a == nil {
		return false
	}

	outAction := h.GetStateAction(out)
	if outAction == nil {
		return false
	}

	a.SetDuration(d, in)
	return h.fsm.AddTransition(id, in, out)
}

func (h *Handler) SetStateDuration(id int, in int, d float32) {
	a := h.GetStateAction(id)
	if a == nil || !h.fsm.CanTransitState(id, in) {
		return
	}
	a.SetDuration(d, in)
}

func (h *Handler) IsValidStateID(id int) bool {
	return id >= 0 && id < len(h.actionList)
}

func (h *Handler) SetCurrentStateAction(id int) bool {
	if !h.IsValidStateID(id) || h.actionList[id] == nil {
		return false
	}
	h.curAction = h.actionList[id]
	return true
}

func (h *Handler) GetCurrentStateId() int {
	return h.fsm.getCurrentStateId()
}

func (h *Handler) GetCurrentStateAction(id int) ActionBase {
	return h.GetStateAction(h.GetCurrentStateId())
}

func (h *Handler) Start(id int) {
	if !h.SetCurrentStateAction(id) {
		return
	}
	h.fsm.SetCurrentState(id)
	h.curAction.OnEnter()

	h.update()
}

func (h *Handler) update() {
	wg := sync.WaitGroup{}

	tick := time.NewTicker(time.Second)

	wg.Add(1)

	go func() {
		defer func() {
			close(h.Ch)
			wg.Done()
		}()

		for {
			select {
			case id := <-h.Ch:
				if id != 100 {
					h.TransitState(id)
				} else {
					return
				}
			case <-tick.C:
				h.UpdateState(1)
			}
		}
	}()
}

func (h *Handler) FindStateAction(in int) ActionBase {
	out := h.fsm.FindOutputStateId(in)
	return h.GetStateAction(out)
}

func (h *Handler) TransitState(in int) bool {
	if h.curAction == nil {
		return false
	}

	old := h.GetCurrentStateId()
	if !h.fsm.TransitState(in) {
		return false
	}

	if old != h.GetCurrentStateId() {
		h.curAction.Expire()
		h.curAction.OnExit()
		h.SetCurrentStateAction(h.GetCurrentStateId())
		h.curAction.OnEnter()
	} else {
		h.curAction.Expire()
	}
	return true
}

func (h *Handler) UpdateState(delta float32) {
	if h.curAction == nil {
		return
	}

	h.curAction.Update(delta)
	if h.curAction.IsExpire() {
		h.curAction.Expire()
		h.TransitState(h.curAction.GetInput())
		return
	}

	h.curAction.OnUpdate(delta)
}

func (h *Handler) HasState(id int) bool {
	if !h.IsValidStateID(id) {
		return false
	}
	return h.actionList[id] != nil
}

func (h *Handler) GetStateCount() int {
	return len(h.actionList)
}

func (h *Handler) AddState(id int, del Delegate) bool {
	if !h.IsValidStateID(id) {
		return false
	}
	if h.actionList[id] != nil {
		return false
	}

	h.actionList[id] = NewAction(del)
	h.fsm.AddState(id)
	return true
}
