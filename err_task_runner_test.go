package kytsya

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRunTasksOk(t *testing.T) {
	var counter uint32

	errCh := NewErrorBox[struct{}]().
		AddTask(func() Result[struct{}] {
			atomic.AddUint32(&counter, 1)

			return Result[struct{}]{}
		}).
		AddTask(func() Result[struct{}] {
			atomic.AddUint32(&counter, 1)

			return Result[struct{}]{}
		}).
		AddTask(func() Result[struct{}] {
			atomic.AddUint32(&counter, 1)

			return Result[struct{}]{}
		}).Run()

	ForChan(errCh, func(val Result[struct{}]) {
		if val.Err != nil {
			t.Fatal(val.Err.Error())
		}

	})

	if atomic.LoadUint32(&counter) != 3 {
		t.Fatal("expect incremented error value")
	}
}

func TestRunTasksError(t *testing.T) {
	var counter uint32

	errCh := NewErrorBox[struct{}]().
		AddTask(func() Result[struct{}] {
			atomic.AddUint32(&counter, 1)

			return Result[struct{}]{}
		}).
		AddTask(func() Result[struct{}] {
			atomic.AddUint32(&counter, 1)

			return Result[struct{}]{}
		}).
		AddTask(func() Result[struct{}] {
			atomic.AddUint32(&counter, 1)

			return Result[struct{}]{
				Err: testError,
			}
		}).Run()

	var expectedErr error

	ForChan(errCh, func(val Result[struct{}]) {
		if val.Err != nil {
			expectedErr = val.Err
		}

	})

	if !errors.Is(expectedErr, testError) {
		t.Fatal("expect test error as result")
	}

	if atomic.LoadUint32(&counter) != 3 {
		t.Fatal("expect incremented error value")
	}
}

func TestRunTasksWithRecover(t *testing.T) {
	res := <-NewErrorBox[struct{}]().WithRecover().AddTask(func() Result[struct{}] {
		panic("🐱")
	}).Run()
	if !errors.Is(res.Err, ErrRecoveredFromPanic) {
		t.Fatal("expect recovered from panic as a result")
	}
}

func BenchmarkErrorBox(b *testing.B) {
	b.Run("pure Go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resCh := make(chan Result[string], 4)
			wg := &sync.WaitGroup{}

			for i := 0; i < 4; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					defer func() {
						if err := recover(); err != nil {
							resCh <- Result[string]{Err: errors.New("error")}
						}
					}()

					resCh <- Result[string]{Data: "🐈"}
				}()
			}

			go func() {
				wg.Wait()
				close(resCh)
			}()

			for range resCh {
			}
		}
	})

	b.Run("kytsunya", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			resCh := NewErrorBox[string]().
				WithRecover().
				AddTask(func() Result[string] {
					return Result[string]{Data: "🐈"}
				}).
				AddTask(func() Result[string] {
					return Result[string]{Data: "🐈"}
				}).
				AddTask(func() Result[string] {
					return Result[string]{Data: "🐈"}
				}).
				AddTask(func() Result[string] {
					return Result[string]{Data: "🐈"}
				}).Run()

			ForChan(resCh, func(val Result[string]) {})
		}
	})
}
