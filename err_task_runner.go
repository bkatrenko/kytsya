package kytsya

import "sync"

type (
	// ErrTaskRunner is a structure for run a group of async tasks that should return a result or error, as
	// run & collect result from some number of network calls.
	// For example:
	// 	resCh := NewErrorBox[string]().
	// 	WithRecover().
	// 	AddTask(func() Result[string] {
	// 		return Result[string]{Data: "1"}
	// 	}).
	// 	AddTask(func() Result[string] {
	// 		return Result[string]{Data: "2"}
	// 	}).
	// 	AddTask(func() Result[string] {
	// 		return Result[string]{Err: errors.New("Houston, we have a problem")}
	// 	}).
	// 	AddTask(func() Result[string] {
	// 		panic("aaaaa")
	// 	}).
	// 	Run()
	// 	// ResCh will be closed after all tasks are done!
	// 	for v := range resCh {
	// 		fmt.Println(v)
	// 	}
	// In this case, WithRecover() returns error to a handler-channel as a Result{Err: val} from function that called
	// a panic.
	ErrTaskRunner[T any] struct {
		tasks  []func() Result[T]
		runner ErrRunner[T]
	}
)

// NewErrorBox is a constructor for task runner.
// We call it *Box cause 🐈🐈🐈 are in love with boxes!
func NewErrorBox[T any]() *ErrTaskRunner[T] {
	return &ErrTaskRunner[T]{
		runner: ErrRunner[T]{
			wg: &sync.WaitGroup{},
		},
	}
}

// WithRecover adds a recovery handler for every task.
// Panic message and a stacktrace will be returned as Result{Err: "error message/panic trace"}
func (tr *ErrTaskRunner[T]) WithRecover() *ErrTaskRunner[T] {
	tr.runner.recover = true
	return tr
}

// AddTask accept a generic function that will contain result or error structure.
func (tr *ErrTaskRunner[T]) AddTask(f func() Result[T]) *ErrTaskRunner[T] {
	tr.tasks = append(tr.tasks, f)
	tr.runner.wg.Add(1)

	return tr
}

// Run spawns all tasks and return a chan Result[T] to collect all.
func (tr *ErrTaskRunner[T]) Run() chan Result[T] {
	tr.runner.errCh = make(chan Result[T], len(tr.tasks))

	for _, task := range tr.tasks {
		tr.runner.Spawn(task)
	}

	go func() {
		tr.runner.wg.Wait()
		close(tr.runner.errCh)
	}()

	return tr.runner.errCh
}
