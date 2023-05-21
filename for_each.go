package kytsya

func ForChan[T any](ch chan T, f func(val T)) {
	for val := range ch {
		f(val)
	}
}

func ForErrChan[T any](ch chan T, f func(val T) error) error {
	for val := range ch {
		err := f(val)
		if err != nil {
			return err
		}
	}

	return nil
}

func ForEach[T any](data []T, f func(i int, val T)) {
	for i, val := range data {
		f(i, val)
	}
}

func ForEachErr[T any](data []T, f func(i int, val T) error) error {
	for i, val := range data {
		if err := f(i, val); err != nil {
			return err
		}
	}

	return nil
}

func Map[T any, V any](data []T, f func(i int, val T) V) []V {
	res := make([]V, 0, len(data))

	for i, val := range data {
		res = append(res, f(i, val))
	}

	return res
}

func MapErr[T any, V any](data []T, f func(i int, val T) (V, error)) ([]V, error) {
	res := make([]V, 0, len(data))

	for i, val := range data {
		mapped, err := f(i, val)
		if err != nil {
			return res, err
		}

		res = append(res, mapped)
	}

	return res, nil
}

func Filter[T any](input []T, f func(i int, val T) bool) []T {
	res := make([]T, 0, len(input))
	for i, val := range input {
		include := f(i, val)

		if include {
			res = append(res, val)
		}
	}

	return res
}

func Reduce[T any, V any](input []T, f func(val T, acc V) V) V {
	var acc V

	for _, val := range input {
		acc = f(val, acc)
	}

	return acc
}
