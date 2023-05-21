package kytsya

import (
	"fmt"
	"runtime/debug"
	"sync"
)

type (
	// Runner could be used to start a new goroutine with ability to recover and print stack trace
	// in case of panic or use a wait group to wait until goroutine function will return.
	// Could be used as a functional/graceful substitution for repeated:

	// defer func(){
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}
	// }()
	//
	// With kytsya it's possible to generalize such operation in one defined manner:
	// kytsya.Goroutine().WithRecover().Spawn(func() {fmt.Println("🐈🐈🐈")}).Wait()
	// Wait group also built-in.
	Runner struct {
		inBox         bool
		shouldRecover bool

		wg *sync.WaitGroup
	}
)

// Goroutine creates a new goroutine runner.
// Example:
// kytsunya2.Goroutine().WithRecover().WithWaitGroup().Spawn(func() {fmt.Println("🐈🐈🐈")}).Wait()
func Goroutine() *Runner {
	return &Runner{}
}

// WithRecover will add defer recover() function to the executor that will recover panics and will
// print the stack trace into stdout.
func (r *Runner) WithRecover() *Runner {
	r.shouldRecover = true

	return r
}

// WithWaitGroup add a wait group into the executor (Wait() could be called, and will block until
// created goroutine will return).
func (r *Runner) WithWaitGroup() *Runner {
	r.wg = &sync.WaitGroup{}

	return r
}

// Spawn start a new goroutine, accept function that will be executed in the newly created routine.
func (r *Runner) Spawn(f func()) Waiter {
	if r.wg != nil && !r.inBox {
		r.wg.Add(1)
	}

	go func() {
		if r.wg != nil {
			defer r.wg.Done()
		}

		if r.shouldRecover {
			defer recoverPrintIfPanic()
		}

		f()
	}()

	return r.wg
}

// RecoverPrintIfPanic recover in case of goroutine starting panic and log the error message with a
// stacktrace.
func recoverPrintIfPanic() {
	if err := recover(); err != nil {
		fmt.Println(fmt.Errorf("%w: %s from %s",
			ErrRecoveredFromPanic,
			err,
			string(debug.Stack()), // debug.Stack is actually version of runtime.Stack that returns formatted output
		))
	}
}
