package fsm

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// StateFunc is a function that runs a given state and returns the next state
type StateFunc func() (StateFunc, error)

type LogFunc func(fmt string, args ...interface{})

type FsmRunner interface {
	// Run starts the execution of the state machine. It returns when a StateFunc returns nil or an error.
	Run() error

	// States is a history of executed states
	States() []string

	LoggingEnabled(on bool)
}

// fsmRunner implements the FsmRunner interface
type fsmRunner struct {
	// name is the name of the state machine
	name string

	// startFunc is the function that starts the state machine
	startFunc StateFunc

	// currentFunc is the StateFunc that is currently running
	currentFunc StateFunc

	// states is a history of executed states
	states []string

	// resetFunc is a function that resets the state machine
	resetFunc func()

	// loggingEnabled is a flag that indicates whether logging is enabled
	loggingEnabled bool

	// logFunc is the function that is called to log messages
	logger LogFunc

	sync.Mutex
}

// Option enables optional arguments for a new runner
type Option func(r *fsmRunner)

func Reset(f func()) Option {
	return func(r *fsmRunner) {
		r.resetFunc = f
	}
}

func Logging(f LogFunc) Option {
	return func(r *fsmRunner) {
		r.logger = f
	}
}

func New(name string, startFunc StateFunc, options ...Option) FsmRunner {
	r := &fsmRunner{
		name:      name,
		startFunc: startFunc,
	}
	for _, option := range options {
		option(r)
	}
	return r
}

func (r *fsmRunner) reset() {
	r.Lock()
	defer r.Unlock()
	r.states = make([]string, 0)
	r.currentFunc = nil
	if r.resetFunc != nil {
		r.resetFunc()
	}
}

func (r *fsmRunner) Run() error {
	defer func() {
		if r.loggingEnabled {
			r.log("State call history:")
			for _, s := range r.states {
				r.log("\t%s", s)
			}
		}
	}()

	r.reset()
	f := r.startFunc
	var err error

	for {
		f, err = r.funcWrapper(f)

		switch {
		case f == nil && err == nil:
			r.log("Run() completed with no errors")
			return nil
		case err != nil:
			r.log("Run() completed with error: %q", err)
			return err
		}

	}
}

// funcWrapper emits logs and updates the state history before running a state
func (r *fsmRunner) funcWrapper(f StateFunc) (StateFunc, error) {
	r.Lock()
	name := getFuncName(f)
	r.states = append(r.states, name)
	r.currentFunc = f
	r.Unlock()

	r.logger("StateFunc %s starting", name)

	f, err := f()

	r.logger("StateFunc %s completed", name)
	return f, err
}

func (r *fsmRunner) States() []string {
	r.Lock()
	defer r.Unlock()
	return r.states
}

func (r *fsmRunner) LoggingEnabled(b bool) {
	if r.logger == nil {
		return
	}

	r.Lock()
	defer r.Unlock()
	r.loggingEnabled = b
}

func (r *fsmRunner) log(s string, i ...interface{}) {
	if r.loggingEnabled && r.logger != nil {
		r.logger(fmt.Sprintf("FSM[%s]: %s", r.name, s), i...)
	}
}

func getFuncName(f StateFunc) string {
	v := reflect.ValueOf(f)
	pc := runtime.FuncForPC(v.Pointer())
	sp := strings.SplitAfter(pc.Name(), ".")
	return strings.TrimSuffix(sp[len(sp)-1], "-fm")
}
