package kytsya

import "errors"

type (
	// TaskRunner represents an executor for a group of goroutines.
	// Could be useful in case of needing to run a controlled group of goroutines with one handler point,
	// panic-handlers of wait group. Constructor function is NewBox 'cause boxes is best friends of cats!
	// For example:
	//
	// NewBox().WithWaitGroup().
	// AddTask(func() {}).
	// AddTask(func() {}).
	// AddTask(func() {}).
	// AfterAll(func() {}).
	// Run().Wait()
	//
	// WithWaitGroup().WithRecover() and AfterAll(f) is not necessary calls, added here just to show the possibilities.
	TaskRunner struct {
		tasks     []func()
		afterFunc func()

		runner Runner
	}

	// Waiter/RunnerFunc/TaskFiller is a group of interfaces that prevents misuse of TaskRunner.
	Waiter interface {
		Wait()
	}

	RunnerFunc interface {
		Run() Waiter
	}

	TaskFiller interface {
		AddTask(f func()) TaskFiller
		Run() Waiter
		AfterAll(f func()) RunnerFunc
	}
)

var (
	ErrWaitWithoutWaitGroup = errors.New("Wait() was called without WithWaitGroup initialization, please, call WithWaitGroup() before")
)

// NewBox creates new runner that controls a set of running goroutines that returns no values.
// It has functionality to:
// 1. Recover all panics with WithRecover()
// 2. Add Wait group to run with WithWaitGroup()
// 3. Add new task to execution with AddTask(f())
// 4. Run the set with Run()
// 5. Wait till all goroutines are done with Wait()
// 6. Add a function that will be executed after all goroutines will done with AfterAll(f())
func NewBox() *TaskRunner {
	return &TaskRunner{
		runner: Runner{
			inBox: true,
		},
	}
}

// WithRecover add a recovery handler to a task funner. Handler will be assigned to every goroutine and in case
// of panic will recover and print a stacktrace into stdout.
func (tr *TaskRunner) WithRecover() *TaskRunner {
	tr.runner.shouldRecover = true

	return tr
}

// WithWaitGroup add a WaitGroup to an executions and makes possible to call Wait() to wait until all tasks will
// done.
func (tr *TaskRunner) WithWaitGroup() *TaskRunner {
	tr.runner.WithWaitGroup()

	return tr
}

// AddTask accept a new task for async execution.
func (tr *TaskRunner) AddTask(f func()) TaskFiller {
	tr.tasks = append(tr.tasks, f)

	if tr.runner.wg != nil {
		tr.runner.wg.Add(1)
	}

	return tr
}

// Run spawns all tasks in a loop.
func (tr *TaskRunner) Run() Waiter {
	for _, task := range tr.tasks {
		tr.runner.Spawn(task)
	}

	return tr
}

// Wait blocks until all tasks will done, call a panic in case of "WithWaitGroup()" was no called.
func (tr *TaskRunner) Wait() {
	if tr.runner.wg == nil {
		panic(ErrWaitWithoutWaitGroup)
	}

	tr.runner.wg.Wait()
}

// AfterAll accept a handler that will be executed after all tasks. In general case it could be used for a
// range of tasks:
// - Close any kinds of connections
// - Logs the results of measure the time of executions
// - close channels
// - etc. up to user.
//
// WARNING:
// The nature of AfterAll is async. It means, that it will wait unit WaitGroup unblock and will be
// executed asynchronously. If it is necessary to wait until AfterFunc will exit, use any possible sync mechanism.
func (tr *TaskRunner) AfterAll(f func()) RunnerFunc {
	tr.afterFunc = f
	tr.WithWaitGroup()
	tr.runner.wg.Add(len(tr.tasks))

	go func() {
		tr.runner.wg.Wait()
		tr.afterFunc()
	}()

	return tr
}
