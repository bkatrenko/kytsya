package kytsya

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"testing"
)

var (
	errTestErrChan = errors.New("error while handle messages from channel")
)

func TestForChan(t *testing.T) {
	dataCh := make(chan struct{})

	NewBox().
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AfterAll(func() {
			close(dataCh)
		}).Run()

	var counter uint32

	ForChan(dataCh, func(val struct{}) {
		counter++
	})

	if counter != 3 {
		t.Fatal("expect 3 struct{}{}s")
	}
}

func TestForChanClosedChan(t *testing.T) {
	dataCh := make(chan struct{})
	close(dataCh)

	ForChan(dataCh, func(val struct{}) {
		t.Fatal("expect no messages from closed channel")
	})
}

func TestForChanClosedBufferedChan(t *testing.T) {
	dataCh := make(chan struct{}, 1)

	dataCh <- struct{}{}
	close(dataCh)

	var counter uint32

	ForChan(dataCh, func(val struct{}) {
		counter++
	})

	if counter != 1 {
		t.Fatal("expect 1 struct{}{}s")
	}
}

func TestForErrChanError(t *testing.T) {
	dataCh := make(chan struct{})

	NewBox().
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AfterAll(func() {
			close(dataCh)
		}).Run()

	var counter uint32

	err := ForErrChan(dataCh, func(val struct{}) error {
		counter++

		if counter == 2 {
			return errTestErrChan
		}

		return nil
	})

	if counter != 2 {
		t.Fatal("expect 2 struct{}{}s")
	}

	if !errors.Is(err, errTestErrChan) {
		t.Fatal("expect errTestErrChan as an output")
	}
}

func TestForErrChanOK(t *testing.T) {
	dataCh := make(chan struct{})

	NewBox().
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AddTask(func() {
			dataCh <- struct{}{}
		}).
		AfterAll(func() {
			close(dataCh)
		}).
		Run()

	var counter uint32

	err := ForErrChan(dataCh, func(val struct{}) error {
		counter++

		return nil
	})

	if counter != 2 {
		t.Fatal("expect 2 struct{}{}s")
	}

	if err != nil {
		t.Fatal("expect no errors")
	}
}

func TestForEach(t *testing.T) {
	data := make([]uint32, 10, 10)

	ForEach(data, func(i int, val uint32) {
		data[i] = val + 1
	})

	ForEach(data, func(i int, val uint32) {
		if val != 1 {
			t.Fatal("expect all values to be \"1\"")
		}
	})
}

func TestForEachErr_Error(t *testing.T) {
	sl := []uint32{1, 2, 3, 4, 5}

	err := ForEachErr(sl, func(i int, val uint32) error {
		if val == 5 {
			return errors.New("bad cat!")
		}

		sl[i]++

		return nil
	})

	if err == nil {
		t.Fatal("expect error")
	}

	expectation := []uint32{2, 3, 4, 5, 5}
	if !reflect.DeepEqual(sl, expectation) {
		t.Fatalf("result is not expected: \ngot: %v\nwant: %v\n", sl, expectation)
	}
}

func TestForEachErr_OK(t *testing.T) {
	sl := []uint32{1, 2, 3, 4, 5}

	err := ForEachErr(sl, func(i int, val uint32) error {
		sl[i]++

		return nil
	})

	if err != nil {
		t.Fatal("expect no errors")
	}

	expectation := []uint32{2, 3, 4, 5, 6}
	if !reflect.DeepEqual(sl, expectation) {
		t.Fatalf("result is not expected: \ngot: %v\nwant: %v\n", sl, expectation)
	}
}

func TestMap(t *testing.T) {
	input := []uint32{1, 2, 3, 4, 5}
	out := Map(input, func(i int, val uint32) string {
		return fmt.Sprint(val)
	})

	expectation := []string{"1", "2", "3", "4", "5"}
	if !reflect.DeepEqual(out, expectation) {
		t.Fatalf("result is not expected: \ngot: %v\nwant: %v\n", out, expectation)
	}
}

func TestMapErr_Error(t *testing.T) {
	input := []uint32{1, 2, 3, 4, 5}
	out, err := MapErr(input, func(i int, val uint32) (string, error) {
		if val == 2 {
			return "", errors.New("bad cat!")
		}

		return fmt.Sprint(val), nil
	})

	if err == nil {
		t.Fatal("expect an error")
	}

	if out[0] != "1" {
		t.Fatal("expect first value modified")
	}
}

func TestMapErr_OK(t *testing.T) {
	input := []uint32{1, 2, 3, 4, 5}
	out, err := MapErr(input, func(i int, val uint32) (string, error) {
		return fmt.Sprint(val), nil
	})

	if err != nil {
		t.Fatal("expect no errors")
	}

	expectation := []string{"1", "2", "3", "4", "5"}
	if !reflect.DeepEqual(out, expectation) {
		t.Fatalf("result is not expected: \ngot: %v\nwant: %v\n", out, expectation)
	}
}

func TestFilter(t *testing.T) {
	out := Filter(
		[]int32{-2, -1, 0, 1, 2},
		func(i int, val int32) bool {
			if val < 0 {
				return false
			}

			return true
		})

	expectation := []int32{0, 1, 2}
	if !reflect.DeepEqual(out, expectation) {
		t.Fatalf("result is not expected: \ngot: %v\nwant: %v\n", out, expectation)
	}
}

func TestReduce(t *testing.T) {
	out := Reduce([]uint32{1, 2, 3, 4, 5}, func(val uint32, acc uint32) uint32 {
		return val + acc
	})

	expectation := 1 + 2 + 3 + 4 + 5
	if out != uint32(expectation) {
		t.Fatalf("result is not expected: \ngot: %v\nwant: %v\n", out, expectation)
	}
}

func TestChanToSlice(t *testing.T) {
	ch := NewErrorBox[int]().
		AddTask(func() Result[int] {
			return Result[int]{Data: 1}
		}).
		AddTask(func() Result[int] {
			return Result[int]{Data: 2}
		}).
		AddTask(func() Result[int] {
			return Result[int]{Data: 3}
		}).
		AddTask(func() Result[int] {
			return Result[int]{Data: 4}
		}).
		AddTask(func() Result[int] {
			return Result[int]{Data: 5}
		}).
		Run()

	out := ChanToSlice[Result[int]](ch)
	ints := Map[Result[int], int](out, func(i int, val Result[int]) int {
		return val.Data
	})

	expectation := []int{1, 2, 3, 4, 5}

	sort.Slice(ints, func(i, j int) bool { return ints[i] < ints[j] })

	if !reflect.DeepEqual(ints, expectation) {
		t.Fatalf("result is not expected: \ngot: %v\nwant: %v\n", ints, expectation)
	}
}
