package kytsya

import "errors"

var (
	ErrRecoveredFromPanic = errors.New("kytsunya: recovered from panic")
	ErrTimeout            = errors.New("kytsunya: goroutine timed out")
)
