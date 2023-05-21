package kytsya

import (
	"fmt"
	"runtime/debug"
	"sync"
	"time"
)

type (
	// ErrRunner is a backend for ErrTaskRunner.
	// Contains all functionality that kytsya needs for running a group of goroutines with defined functionality.
	// Could be used also separately as Erroutine (aka goroutine that returns a result of error).
	// For example:
	//
	// res := Erroutine[string]().
	// 	WithRecover().
	// 	WithTimeout(time.Second).
	// 	Spawn(
	// 		func() Result[string] {
	// 			time.Sleep(time.Second * 10)
	// 			return Result[string]{Data: "🐈"}
	// 		}).
	// 	WaitAsync()

	//	fmt.Println(<-res)
	//  In this case goroutine also using WithTimeout functionality.
	ErrRunner[T any] struct {
		recover bool

		timeout time.Duration
		wg      *sync.WaitGroup

		errCh chan Result[T]
	}

	// Result[T] define kind of OneOf:
	// Generic result of the execution or error message.
	// Error message will be also returned in case of panic or timeout.
	Result[T any] struct {
		Data T
		Err  error
	}
)

func Erroutine[T any]() *ErrRunner[T] {
	return &ErrRunner[T]{
		errCh: make(chan Result[T], 1),
	}
}

func (r *ErrRunner[T]) WithRecover() *ErrRunner[T] {
	r.recover = true

	return r
}

func (r *ErrRunner[T]) WithTimeout(timeout time.Duration) *ErrRunner[T] {
	r.timeout = timeout

	return r
}

func (r *ErrRunner[T]) Spawn(f func() Result[T]) *ErrRunner[T] {
	go func() {
		if r.wg != nil {
			defer r.wg.Done()
		}

		if r.recover {
			defer errorIfPanic(r.errCh)
		}

		res := f()

		r.errCh <- res
	}()

	return r
}

func (r *ErrRunner[T]) Wait() Result[T] {
	if r.timeout == 0 {
		return <-r.errCh
	}

	select {
	case err := <-r.errCh:
		return err
	case <-time.After(r.timeout):
		return Result[T]{
			Err: ErrTimeout,
		}
	}
}

func (r *ErrRunner[T]) WaitAsync() chan Result[T] {
	resCh := make(chan Result[T], 1)

	go func() {
		defer close(resCh)

		if r.timeout == 0 {
			resCh <- <-r.errCh
		}

		select {
		case res := <-r.errCh:
			resCh <- res
		case <-time.After(r.timeout):
			resCh <- Result[T]{
				Err: ErrTimeout,
			}
		}
	}()

	return resCh
}

func errorIfPanic[T any](errCh chan Result[T]) {
	if err := recover(); err != nil {
		errCh <- Result[T]{
			Err: fmt.Errorf("%w: %s from %s",
				ErrRecoveredFromPanic,
				err,
				string(debug.Stack()), // debug.Stack is actually version of runtime.Stack that returns formatted output
			),
		}
	}
}
