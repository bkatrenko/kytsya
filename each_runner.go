package kytsya

type (
	EachRunner[T any, V any] struct {
		data    []T
		recover bool
	}
)

func NewEachRunner[T any, V any](data []T) *EachRunner[T, V] {
	return &EachRunner[T, V]{
		data: data,
	}
}

func (er *EachRunner[T, V]) WithRecover() *EachRunner[T, V] {
	er.recover = true

	return er
}

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
