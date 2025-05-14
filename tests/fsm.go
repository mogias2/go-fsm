package main

import (
	"fmt"
	"sync"

	"github.com/mogias2/go-fsm/internal/fsm"
)

type RunDelegate struct {
	f *FSMTest
}

func (r *RunDelegate) Enter() {
	fmt.Println("enter run")
}

func (r *RunDelegate) Update(delta float32) {
	fmt.Println("update run")
}

func (r *RunDelegate) Exit() {
	r.f.ch <- 1
	fmt.Println("exit run")
}

type RunAction struct {
	fsm.Action
}

var _ fsm.ActionBase = (*RunAction)(nil)

func (r *RunAction) OnUpdate(delta float32) {
	if r.Del == nil {
		return
	}
	r.Del.Update(delta)
}

type FSMTest struct {
	handler *fsm.Handler
	ch      chan int
	wg      sync.WaitGroup
}

func (f *FSMTest) init() {
	f.ch = make(chan int)

	f.handler = fsm.NewHandler(3, "test")

	idle := &fsm.Action{Name: "idle"}

	run := &RunAction{}
	run.Name = "run"
	run.Del = &RunDelegate{f: f}

	fly := &fsm.Action{Name: "fly"}

	f.handler.AddStateWithAction(0, idle)
	f.handler.AddStateWithAction(1, run)
	f.handler.AddStateWithAction(2, fly)

	f.handler.AddTransition(0, 1, 1, 0)
	f.handler.AddTransition(1, 1, 2, 3)

	f.handler.Start(0)
}

func (f *FSMTest) start() {
	f.wg = sync.WaitGroup{}

	f.wg.Add(1)
	go func() {
		defer f.wg.Done()

		f.PrintAction()

		f.handler.Ch <- 1

		select {
		case _ = <-f.ch:
			f.PrintAction()
			f.handler.Ch <- 100
			break
		}
	}()

	f.wg.Wait()

	close(f.ch)
}

func (f *FSMTest) PrintAction() {
	id := f.handler.GetCurrentStateId()
	if a, ok := f.handler.GetCurrentStateAction(id).(*fsm.Action); ok {
		fmt.Println(a.Name)
	}
}

func (f *FSMTest) GetCurrentState() int {
	return f.handler.GetCurrentStateId()
}
