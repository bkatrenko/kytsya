package kytsya

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBoxRun(t *testing.T) {
	var (
		counter uint32
		wg      *sync.WaitGroup
	)

	wg = &sync.WaitGroup{}
	wg.Add(5)

	NewBox().
		AddTask(func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}).
		AddTask(func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}).
		AddTask(func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}).
		AddTask(func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}).
		AddTask(func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}).Run()

	wg.Wait()

	if atomic.LoadUint32(&counter) != 5 {
		t.Fatal("expected counter to be 5")
	}
}

func TestBoxRunRecover(t *testing.T) {
	var (
		counter uint32
		wg      *sync.WaitGroup
	)

	wg = &sync.WaitGroup{}
	wg.Add(2)

	NewBox().WithRecover().
		AddTask(func() {
			defer wg.Done()
			atomic.AddUint32(&counter, 1)
		}).
		AddTask(func() {
			defer wg.Done()
			panic("aaa")
		}).Run()

	wg.Wait()

	if atomic.LoadUint32(&counter) != 1 {
		t.Fatal("expected counter to be 1")
	}
}

func TestBoxRunWaitRecover(t *testing.T) {
	var (
		counter uint32
	)

	NewBox().
		WithWaitGroup().
		WithRecover().
		AddTask(func() {
			atomic.AddUint32(&counter, 1)
		}).
		AddTask(func() {
			panic("aaa")
		}).Run().Wait()

	if atomic.LoadUint32(&counter) != 1 {
		t.Fatal("expected counter to be 1")
	}
}

func TestBoxRunWait(t *testing.T) {
	var (
		counter uint32
	)

	box := NewBox().WithWaitGroup()

	for i := 0; i < 5; i++ {
		box.AddTask(func() {
			atomic.AddUint32(&counter, 1)
		})
	}

	box.Run().Wait()

	if atomic.LoadUint32(&counter) != 5 {
		t.Fatal("expected counter to be 5")
	}
}

func TestAfterAll(t *testing.T) {
	type timeSaver struct {
		number uint32
		ti     time.Time
	}

	var (
		counter    uint32
		timeArray  = []timeSaver{{}, {}, {}, {}}
		afterAllCh = make(chan struct{})
	)

	NewBox().
		AddTask(func() {
			atomic.AddUint32(&counter, 1)
			timeArray[0] = timeSaver{0, time.Now()}
		}).
		AddTask(func() {
			atomic.AddUint32(&counter, 1)
			timeArray[1] = timeSaver{1, time.Now()}
		}).
		AddTask(func() {
			atomic.AddUint32(&counter, 1)
			timeArray[2] = timeSaver{2, time.Now()}
		}).
		AfterAll(func() {
			atomic.AddUint32(&counter, 1)
			timeArray[3] = timeSaver{3, time.Now()}
			afterAllCh <- struct{}{}
		}).
		Run().
		Wait()

	<-afterAllCh
	sort.Slice(timeArray, func(i, j int) bool {
		return timeArray[i].ti.UnixNano() < timeArray[j].ti.UnixNano()
	})

	if timeArray[3].number != 3 {
		fmt.Println(timeArray)
		t.Fatal("after all number expected to be the last one")
	}
}

func TestRunnerMisuse_NilWaitGroup(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			if !errors.Is(err.(error), ErrWaitWithoutWaitGroup) {
				t.Fatal("expect panic 'cause WaitGroup was not initialized")
			}
		}
	}()

	NewBox().AddTask(func() {}).Run().Wait()
}

func BenchmarkTestNewBoxWithWait(b *testing.B) {
	b.Run("pure Go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var (
				counter uint32
			)

			wg := &sync.WaitGroup{}

			for i := 0; i < 100; i++ {
				wg.Add(1)

				go func() {
					defer wg.Done()

					s := ""
					s += "salt"
					if strings.Contains(s, "s") {
						atomic.AddUint32(&counter, 1)
					}

				}()
			}

			wg.Wait()
		}
	})

	b.Run("kytsunya", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var (
				counter uint32
			)

			box := NewBox().WithWaitGroup()

			for i := 0; i < 100; i++ {
				box.AddTask(func() {
					s := ""
					s += "salt"
					if strings.Contains(s, "s") {
						atomic.AddUint32(&counter, 1)
					}
				})
			}

			box.Run().Wait()
		}
	})
}
