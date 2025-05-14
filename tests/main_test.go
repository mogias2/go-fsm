package main

import (
	"testing"
)

func TestFSM(t *testing.T) {
	fsm := &FSMTest{}
	fsm.init()
	fsm.start()

	if fsm.GetCurrentState() != 2 {
		t.Error("wrong value")
	}
}
