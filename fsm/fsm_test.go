package fsm

import (
	"fmt"
	"runtime"
	"testing"
)

type SM struct {
	err    bool
	cTrace []string
}

func (sm *SM) Start() (StateFunc, error) {

	sm.trace()
	return sm.Middle, nil
}

func (sm *SM) Middle() (StateFunc, error) {
	sm.trace()
	return sm.End, nil
}

func (sm *SM) End() (StateFunc, error) {
	sm.trace()
	return nil, nil
}

func (sm *SM) reset() {
	sm.cTrace = nil
}

func (sm *SM) trace() {
	pc, _, _, _ := runtime.Caller(1)
	sm.cTrace = append(sm.cTrace, runtime.FuncForPC(pc).Name())
}

type logging struct {
	msgs []string
}

func (l *logging) Log(s string, i ...interface{}) {
	l.msgs = append(l.msgs, fmt.Sprintf(s, i...))
}

func TestRunner(t *testing.T) {
	sm := &SM{}
	l := &logging{}
	exec := New("tester", sm.Start, Reset(sm.reset), Logging(l.Log))
	exec.Run()
}
