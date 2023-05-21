package kytsya

import (
	"reflect"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
)

func TestRunEach(t *testing.T) {
	resCh := NewEachRunner[int, int]([]int{1, 2, 3, 4}).Handle(func(val int) Result[int] {
		return Result[int]{
			Data: val + 1,
		}
	})

	result := []int{}
	expected := []int{2, 3, 4, 5}

	for res := range resCh {
		result = append(result, res.Data)
	}

	sort.Ints(result)

	if !reflect.DeepEqual(expected, []int{2, 3, 4, 5}) {
		t.Fatal("expected values are:", expected, "got:", result)
	}
}

func TestEmpty(t *testing.T) {
	resCh := NewEachRunner[int, int](nil).Handle(func(val int) Result[int] {
		return Result[int]{
			Data: val + 1,
		}
	})

	result := []int{}

	for res := range resCh {
		result = append(result, res.Data)
	}

	if len(result) != 0 {
		t.Fatal("expect the empty list")
	}
}

func TestEachRunnerWithRecover(t *testing.T) {
	res := NewEachRunner[int, int]([]int{1, 2, 3}).WithRecover().Handle(func(val int) Result[int] {
		if val == 1 {
			panic("aaa")
		}

		return Result[int]{Data: val + 1}
	})

	ForChan(res, func(val Result[int]) {
		if val.Data == 3 || val.Data == 4 {
			if val.Err != nil {
				t.Fatal("expect no error")
			}

			return
		}

		if val.Err == nil {
			t.Fatal("expect an error")
		}
	})
}

func BenchmarkRunEach(b *testing.B) {
	b.Run("pure Go", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

			resultCh := make(chan uint32)
			wg := &sync.WaitGroup{}

			for _, value := range data {
				wg.Add(1)

				go func(value uint32) {
					defer wg.Done()
					atomic.AddUint32(&value, 1)
				}(value)
			}

			go func() {
				wg.Wait()
				close(resultCh)
			}()

			for value := range resultCh {
				atomic.AddUint32(&value, 1)
			}
		}
	})

	b.Run("kytsunya", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

			resCh := NewEachRunner[uint32, uint32](data).
				Handle(func(value uint32) Result[uint32] {
					return Result[uint32]{Data: atomic.AddUint32(&value, 1)}
				})

			ForChan(resCh, func(value Result[uint32]) {
				atomic.AddUint32(&value.Data, 1)
			})
		}
	})
}
