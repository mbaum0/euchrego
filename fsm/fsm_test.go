package fsm

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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
	if sm.err {
		return nil, fmt.Errorf("error")
	} else {
		return sm.End, nil
	}
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

func TestRunnerLoggingEnabled(t *testing.T) {
	sm := &SM{}
	l := &logging{}
	exec := New("tester", sm.Start, Reset(sm.reset), Logger(l.Log))
	exec.EnableLogging(true)

	expectedLog :=
		[]string{
			"FSM[tester]: StateFunc Start starting",
			"FSM[tester]: StateFunc Start completed",
			"FSM[tester]: StateFunc Middle starting",
			"FSM[tester]: StateFunc Middle completed",
			"FSM[tester]: StateFunc End starting",
			"FSM[tester]: StateFunc End completed",
			"FSM[tester]: Run() completed with no errors",
			"FSM[tester]: State call history:",
			"FSM[tester]: \tStart",
			"FSM[tester]: \tMiddle",
			"FSM[tester]: \tEnd",
		}
	exec.Run()

	assert.Equal(t, expectedLog, l.msgs)
}

func TestRunnerLoggingDisabled(t *testing.T) {
	sm := &SM{}
	l := &logging{}
	exec := New("tester", sm.Start, Reset(sm.reset), Logger(l.Log))
	exec.EnableLogging(false)
	exec.Run()
	assert.Empty(t, l.msgs)
}

func TestRunnerNoError(t *testing.T) {
	sm := &SM{err: false}
	l := &logging{}
	exec := New("tester", sm.Start, Reset(sm.reset), Logger(l.Log))
	exec.EnableLogging(true)

	expectedLog :=
		[]string{
			"FSM[tester]: StateFunc Start starting",
			"FSM[tester]: StateFunc Start completed",
			"FSM[tester]: StateFunc Middle starting",
			"FSM[tester]: StateFunc Middle completed",
			"FSM[tester]: StateFunc End starting",
			"FSM[tester]: StateFunc End completed",
			"FSM[tester]: Run() completed with no errors",
			"FSM[tester]: State call history:",
			"FSM[tester]: \tStart",
			"FSM[tester]: \tMiddle",
			"FSM[tester]: \tEnd",
		}
	exec.Run()

	assert.Equal(t, expectedLog, l.msgs)
}

func TestRunnerWithError(t *testing.T) {
	sm := &SM{err: true}
	l := &logging{}
	exec := New("tester", sm.Start, Reset(sm.reset), Logger(l.Log))
	exec.EnableLogging(true)

	expectedLog :=
		[]string{
			"FSM[tester]: StateFunc Start starting",
			"FSM[tester]: StateFunc Start completed",
			"FSM[tester]: StateFunc Middle starting",
			"FSM[tester]: StateFunc Middle completed",
			"FSM[tester]: Run() completed with error: \"error\"",
			"FSM[tester]: State call history:",
			"FSM[tester]: \tStart",
			"FSM[tester]: \tMiddle",
		}
	exec.Run()

	assert.Equal(t, expectedLog, l.msgs)
}
