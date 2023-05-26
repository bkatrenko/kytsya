package kytsya

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
)

// Need to run a goroutines? No problem!
func TestRunGoroutines(t *testing.T) {
	/////////////////////////////////////////////////////////////////
	// Need to run a singe goroutine?
	Goroutine().Spawn(func() { fmt.Println("🐱") })

	/////////////////////////////////////////////////////////////////
	// Add recovery handler? Stack trace will be printed to stdout
	Goroutine().
		WithRecover().
		Spawn(func() { fmt.Println("🐱") })

	/////////////////////////////////////////////////////////////////
	// Wait until it's done?
	Goroutine().
		WithRecover().
		WithWaitGroup().
		Spawn(func() { fmt.Println("🐱") }).Wait()
}

// Need to run goroutine that return a result/content/data? No problems!
func TestRunErroutine(t *testing.T) {
	res := <-Erroutine[string]().Spawn(func() Result[string] {
		return Result[string]{Data: "🐱"}
	}).WaitAsync()

	fmt.Println(res)

	/////////////////////////////////////////////////////////////////
	// There are few goroutines we need to run? Cool!
	box := NewErrorBox[string]()

	ForEach([]int{1, 2, 3}, func(i int, val int) {
		box.AddTask(func() Result[string] {
			return Result[string]{Data: strconv.Itoa(val)}
		})
	})

	ForChan(box.Run(), func(val Result[string]) {
		fmt.Println(val.Data, val.Err)
	})

	/////////////////////////////////////////////////////////////////
	// Also could be "one liner"
	resCh := NewErrorBox[string]().
		AddTask(func() Result[string] {
			return Result[string]{Data: strconv.Itoa(1)}
		}).
		AddTask(func() Result[string] {
			return Result[string]{Data: strconv.Itoa(2)}
		}).
		AddTask(func() Result[string] {
			return Result[string]{Data: strconv.Itoa(3)}
		}).Run()

	ForChan(resCh, func(val Result[string]) {
		fmt.Println(val.Data, val.Err)
	})

	/////////////////////////////////////////////////////////////////
	// Also! Could be an EachRunner!
	resCh = NewEachRunner[int, string]([]int{1, 2, 3}).Handle(func(val int) Result[string] {
		return Result[string]{Data: strconv.Itoa(val)}
	})

	ForChan(resCh, func(val Result[string]) {
		fmt.Println(val.Data, val.Err)
	})

	/////////////////////////////////////////////////////////////////
	// Feel some possibility to have a panic? Easy!
	resCh = NewEachRunner[int, string]([]int{1, 2, 3}).
		WithRecover().
		Handle(func(val int) Result[string] {
			if val == 2 {
				panic("Houston, we have a problem")
			}

			return Result[string]{Data: strconv.Itoa(val)}
		})

	ForChan(resCh, func(val Result[string]) {
		fmt.Println(val.Data, val.Err) //<----- panic stacktrace will be returned as an error!
	})
}

// Needs to make gophers life a bit "functional"?
func TestFunctionalEra(t *testing.T) {
	// Range it!
	ForEach([]int{1, 2, 3, 4, 5, 6}, func(i, val int) {
		fmt.Printf("index: %d value: %d", i, val)
	})

	// Filter it!
	// output: [2 4 6]
	fmt.Println(Filter([]int{1, 2, 3, 4, 5, 6}, func(i, val int) bool {
		return val%2 == 0
	}))

	// Map it!
	resMap := Map([]int{1, 2, 3, 4, 5, 6}, func(i, val int) string {
		return strconv.Itoa(val)
	})

	// output: [1 2 3 4 5 6] as an array of string
	fmt.Println(resMap)

	// Reduce it!
	// output: 21
	fmt.Println(Reduce([]int{1, 2, 3, 4, 5, 6}, func(val, acc int) int {
		return val + acc
	}))
}

func TestMapErr(t *testing.T) {
	res, err := MapErr([]int{1, 2, 3, 4, 5}, func(i, val int) (string, error) {
		if val == 3 {
			return "", errors.New("it's a 3!!!")
		}

		return strconv.Itoa(val), nil
	})
	if err != nil {
		//panic(err) could be here
	}

	fmt.Println(res) // result will print empty slice while we got an error during the iteration.
}
