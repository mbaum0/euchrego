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

// FsmRunner is an interface that defines a state machine runner
type FsmRunner interface {
	// Run starts the execution of the state machine. It returns when a StateFunc returns nil or an error.
	Run() error

	// States returns a history of executed states
	States() []string

	// EnableLogging is used to enable or disable logging
	EnableLogging(on bool)
}

// fsmRunner is an implementation of the FsmRunner interface
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

	// stateChangeNotifier is a channel that is used to notify when the state changes
	stateChangeNotifier chan bool

	sync.Mutex
}

// Options are used to configured a new FsmRunner
type Option func(r *fsmRunner)

// Reset is used to set the Runner.Reset() function
func Reset(f func()) Option {
	return func(r *fsmRunner) {
		r.resetFunc = f
	}
}

// Notifier is used to set a notifier channel for when state changes
func Notifier(c chan bool) Option {
	return func(r *fsmRunner) {
		r.stateChangeNotifier = c
	}
}

// Logger is used to set the Runner.Logger() function
func Logger(f LogFunc) Option {
	return func(r *fsmRunner) {
		r.logger = f
	}
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
		if r.stateChangeNotifier != nil {
			r.stateChangeNotifier <- true
		}

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

	r.log("StateFunc %s starting", name)

	f, err := f()

	r.log("StateFunc %s completed", name)
	return f, err
}

// States returns a list of executed states
func (r *fsmRunner) States() []string {
	r.Lock()
	defer r.Unlock()
	return r.states
}

// EnableLogging is used to toggle logging on or off
func (r *fsmRunner) EnableLogging(b bool) {
	if r.logger == nil {
		return
	}

	r.Lock()
	defer r.Unlock()
	r.loggingEnabled = b
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

// log emits a log message if logging is enabled
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
