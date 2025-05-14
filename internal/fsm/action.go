package fsm

type Delegate interface {
	Enter()
	Update(delta float32)
	Exit()
}

type ActionBase interface {
	OnEnter()
	OnExit()
	OnUpdate(delta float32)
	SetDuration(d float32, in int)
	Expire()
	IsExpire() bool
	Update(delta float32)
	GetInput() int
}

type Action struct {
	deltaSeconds float32
	expiration   float32
	input        int
	Del          Delegate
	Name         string
	ActionBase
}

var _ ActionBase = (*Action)(nil)

func NewAction(del Delegate) *Action {
	return &Action{
		Del: del,
	}
}

func (a *Action) SetDuration(d float32, in int) {
	a.deltaSeconds = 0.0
	a.expiration = d
	a.input = in
}

func (a *Action) IsExpire() bool {
	return a.expiration != 0.0 && a.deltaSeconds >= a.expiration
}

func (a *Action) Expire() {
	a.deltaSeconds = 0.0
}

func (a *Action) Clear() {
	a.deltaSeconds = 0.0
	a.expiration = 0.0
}

func (a *Action) GetInput() int {
	return a.input
}

func (a *Action) Update(delta float32) {
	a.deltaSeconds += delta
}

func (a *Action) OnEnter() {
	if a.Del == nil {
		return
	}
	a.Del.Enter()
}

func (a *Action) OnExit() {
	if a.Del == nil {
		return
	}
	a.Del.Exit()
}

func (a *Action) OnUpdate(delta float32) {
	if a.Del == nil {
		return
	}
	a.Del.Update(delta)
}

func (a *Action) SetName(name string) {
	a.Name = name
}
