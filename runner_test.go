package kytsya

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestRun(t *testing.T) {
	var (
		counter uint32
		doneCh  chan struct{}
	)

	doneCh = make(chan struct{})

	Goroutine().Spawn(func() {
		atomic.AddUint32(&counter, 1)
		doneCh <- struct{}{}
	})

	<-doneCh

	if counter != 1 {
		t.Fatal("expected counter to be 1")
	}
}

func TestWithRecover(t *testing.T) {
	var (
		counter uint32
		doneCh  chan struct{}
	)

	doneCh = make(chan struct{})

	Goroutine().
		WithRecover().
		WithWaitGroup().
		Spawn(func() {
			atomic.AddUint32(&counter, 1)

			defer func() {
				doneCh <- struct{}{}
			}()

			panic("Houston, we have a problem")
		})

	<-doneCh

	if counter != 1 {
		t.Fatal("expected counter to be 1")
	}
}

func TestWithWaitGroup(t *testing.T) {
	var (
		counter uint32
	)

	Goroutine().
		WithWaitGroup().
		Spawn(func() {
			atomic.AddUint32(&counter, 1)
		}).Wait()

	if counter != 1 {
		t.Fatal("expected counter to be 1")
	}
}

func BenchmarkTestSpawn(b *testing.B) {
	b.Run("pure Go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg := &sync.WaitGroup{}
			for i := 0; i < 100; i++ {
				var (
					counter uint32
				)

				wg.Add(1)

				go func() {
					defer wg.Done()
					atomic.AddUint32(&counter, 1)
				}()
			}

			wg.Wait()
		}
	})

	b.Run("kytsunya", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			wg := &sync.WaitGroup{}
			for i := 0; i < 100; i++ {
				var (
					counter uint32
				)

				wg.Add(1)

				Goroutine().Spawn(func() {
					defer wg.Done()
					atomic.AddUint32(&counter, 1)
				})
			}

			wg.Wait()
		}
	})
}
