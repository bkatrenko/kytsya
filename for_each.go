package kytsya

// ForChan will range through the channel until its closed
func ForChan[T any](ch chan T, f func(val T)) {
	for val := range ch {
		f(val)
	}
}

// ForErrChan will range through the channel entries and handler will receive every value.
// In case of error returned from the handler, loop will be stopped.
func ForErrChan[T any](ch chan T, f func(val T) error) error {
	for val := range ch {
		err := f(val)
		if err != nil {
			return err
		}
	}

	return nil
}

// ForEach do a range through the slice of data, and will execute input handler function on every
// member of the input slice.
// For example:
//
//	ForEach([]int{1, 2, 3, 4, 5, 6}, func(i, val int) {
//		fmt.Printf("index: %d value: %d", i, val)
//	})
//
// Will print index and value of every member of input slice.
func ForEach[T any](data []T, f func(i int, val T)) {
	for i, val := range data {
		f(i, val)
	}
}

// ForEachErr do a range through the slice of data, and will execute input handler function on every
// member of the input slice.
// For example:
//
// ForEachErr([]int{1, 2, 3, 4, 5, 6}, func(i, val int) error {
// 	if val == 3 {
// 		return errors.New("it's a 3!")
// 	}

//		fmt.Printf("index: %d value: %d", i, val)
//		return nil
//	})
//
// Will print index and value of every member of input slice.
// In case of error, iteration will be interrupted.
func ForEachErr[T any](data []T, f func(i int, val T) error) error {
	for i, val := range data {
		if err := f(i, val); err != nil {
			return err
		}
	}

	return nil
}

// Map receive a slice of data of one type, apply handler function to every element of it,
// and return a slice of data of another type (modified by handler).
// For example:

// resMap := Map([]int{1, 2, 3, 4, 5, 6}, func(i, val int) string {
// 	return strconv.Itoa(val)
// })

// output: [1 2 3 4 5 6] as an array of string
// In this example Map modified an array of integers to an array of string.
func Map[T any, V any](data []T, f func(i int, val T) V) []V {
	res := make([]V, 0, len(data))

	for i, val := range data {
		res = append(res, f(i, val))
	}

	return res
}

// MapErr works similar to map with out exception: handler function could return an error. In such case,
// iteration will be stopped, and function will return empty result with an error.
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

// Filter function accept a slice of T(any), and a filter handler that returns a bool value (include?).
// In case of true, value will be included in the output slice, otherwise will be filtered out.
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

// Reduce receive a slice of values, and return a single value as a result.
// For example:
//
//	fmt.Println(Reduce([]int{1, 2, 3, 4, 5, 6}, func(val, acc int) int {
//		return val + acc
//	}))
//
// Will print a sum of all integers in the slice.
func Reduce[T any, V any](input []T, f func(val T, acc V) V) V {
	var acc V

	for _, val := range input {
		acc = f(val, acc)
	}

	return acc
}

// ChanToSlice converts channel content to slice.
// Function expect channel to be closed at the end.
func ChanToSlice[T any](input chan T) (res []T) {
	ForChan[T](input, func(val T) {
		res = append(res, val)
	})

	return
}
