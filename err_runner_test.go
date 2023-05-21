package kytsya

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

var (
	testError = errors.New("houston, we have a problem")
)

func TestErroutineOk(t *testing.T) {
	var counter uint32

	err := Erroutine[struct{}]().Spawn(func() Result[struct{}] {
		atomic.AddUint32(&counter, 1)
		return Result[struct{}]{}
	}).Wait()

	if err.Err != nil {
		t.Fatal("expect nil error")
	}

	if atomic.LoadUint32(&counter) != 1 {
		t.Fatal("expect incremented counter")
	}
}

func TestErroutineError(t *testing.T) {
	var counter uint32

	err := Erroutine[struct{}]().Spawn(func() Result[struct{}] {
		atomic.AddUint32(&counter, 1)
		return Result[struct{}]{
			Err: testError,
		}
	}).Wait()

	if err.Err == nil {
		t.Fatal("expect not nil error")
	}

	if err.Err.Error() != testError.Error() {
		t.Fatal("expect test error in response")
	}

	if atomic.LoadUint32(&counter) != 1 {
		t.Fatal("expect incremented counter")
	}
}

func TestErroutineWithRecover(t *testing.T) {
	err := Erroutine[struct{}]().WithRecover().Spawn(func() Result[struct{}] {
		panic(testError)
	}).Wait()

	if err.Err == nil {
		t.Fatal("expect non nil error")
	}

	if !errors.Is(err.Err, ErrRecoveredFromPanic) {
		t.Fatal("expect recovered test error")
	}
}

func TestErroutineWithTimeout(t *testing.T) {
	err := Erroutine[struct{}]().WithTimeout(time.Second).Spawn(func() Result[struct{}] {
		time.Sleep(time.Second * 2)

		return Result[struct{}]{}
	}).Wait()

	if err.Err == nil {
		t.Fatal("expect non nil error")
	}

	if !errors.Is(err.Err, ErrTimeout) {
		t.Fatal("expect timeout error")
	}
}

func TestWaitAsync(t *testing.T) {
	resCh := Erroutine[uint32]().Spawn(func() Result[uint32] {
		return Result[uint32]{Data: 1}
	}).WaitAsync()

	res := <-resCh
	if res.Data != 1 {
		t.Fatal("expected async data got from channel")
	}
}

func TestWaitAsync_Timeout(t *testing.T) {
	resCh := Erroutine[uint32]().WithTimeout(time.Second).
		Spawn(func() Result[uint32] {
			time.Sleep(time.Second * 2)
			return Result[uint32]{Data: 1}
		}).WaitAsync()

	res := <-resCh

	if !errors.Is(res.Err, ErrTimeout) {
		t.Fatal("goroutine expect timeout")
	}
}

func TestWaitAsync_Timeout_OK(t *testing.T) {
	resCh := Erroutine[uint32]().WithTimeout(time.Second).
		Spawn(func() Result[uint32] {
			time.Sleep(time.Millisecond * 50)
			return Result[uint32]{Data: 1}
		}).WaitAsync()

	res := <-resCh

	if res.Data != 1 {
		t.Fatal("expected async data got from channel")
	}
}

func TestWaitWithErrorWithTimeout(t *testing.T) {
	err := errors.New("something went wrong")
	res := Erroutine[int]().WithTimeout(time.Second).Spawn(func() Result[int] {
		return Result[int]{Err: err}
	}).Wait()

	if !errors.Is(res.Err, err) {
		t.Fatal("expect an error")
	}
}
