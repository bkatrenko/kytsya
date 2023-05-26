package kytsya

type (
	// EachRunner could be used in case of needing to handle every list member in different goroutine.
	EachRunner[T any, V any] struct {
		data    []T
		recover bool
	}
)

// NewEachRunner returns an instance of each runner - entity that handle a slice of data with handler,
// where each member handled in a different goroutine.
// T represents input data type.
// V represents output data type.
func NewEachRunner[T any, V any](data []T) *EachRunner[T, V] {
	return &EachRunner[T, V]{
		data: data,
	}
}

// WithRecover adds recovery handler to each spawned goroutine.
func (er *EachRunner[T, V]) WithRecover() *EachRunner[T, V] {
	er.recover = true

	return er
}

// Handle(f) accepts a functional handler for every member on the input list. Handler will be spawned in a separate
// goroutine and will receive on of input list entry.
func (er *EachRunner[T, V]) Handle(handler func(val T) Result[V]) chan Result[V] {
	box := NewErrorBox[V]()

	if er.recover {
		box.WithRecover()
	}

	for _, val := range er.data {
		cpy := val

		box.AddTask(func() Result[V] {
			return handler(cpy)
		})
	}

	return box.Run()
}
