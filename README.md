# 🇺🇦kytsya🇺🇦
![Logo](https://github.com/bkatrenko/kytsya/blob/master/logo.png)

## Go toolkit!
Kytsya means "kitten" in ukrainian. While cats are the best programmers friends, we decided to call the repo in honour of these beautiful creatures.
It's small, but powerful kit that could be used to change an approach to threading and working with slices.

It contains:
* Number of controllers for goroutines to use them in a functional, save and readable manner
* Controllers for run goroutines that should return a result
* Builders for start a goroutines with a wait group, recovery handler or/and timeout
* ForEach, ForEachChan, Map, Reduce, Filter functions to handle data in a functional and readable manner

Cases where kytsya will shine: 
* Work hardly with goroutines? Kytsya helps to manage that in nice way.
* Needs recovery handlers, but don't want to manually control every case that should be covered with deferred recover? Here you could find __WithRecover()__ builder, that makes recovery graceful!
* Run goroutins that return a results / need to make everything stable and built in one manner? Kytsya will help to do it with __ErrorBox__ or __ForEachRunner__.
* Every developer in a big team use its own threading style? Have a mix of channels, waitgroups, "res" slices? Kytsya gives one approach for everything.
* Work with data? __ForEach__, __Map__, __Reduce__ will definetly what is missing in go's std lib.

Kytsya is not a framework, it's a toolkit that gives an ability, but not force developers to use the one way.
Here is only small part from what it is doing:
1. Needs to run a set of goroutines with recovery handlers and wait group?
```
	kytsya.NewBox().
		WithRecover().
		WithWaitGroup().
		AddTask(func() { fmt.Println("🐈") }).
		AddTask(func() { fmt.Println("🐈") }).
		AddTask(func() { fmt.Println("🐈") }).
		Run().Wait()
```
2. Needs in a safe way to run a set of goroutines and properly read results?
```
	resCh := NewErrorBox[string]().
		WithRecover().
		AddTask(func() Result[string] {
			return Result[string]{Data: "🐈"}
		}).
		AddTask(func() Result[string] {
			panic("dog detected")
			return Result[string]{Data: "🐕"}
		}).Run()
```
Read the output:
```
	ForChan(resCh, func(val Result[string]) {
		fmt.Println(val)
	})
```
That will print:
```
{🐈 <nil>}
{ kytsunya: recovered from panic: dog detected from goroutine 22 [running]:
runtime/debug.Stack()
	/opt/homebrew/Cellar/go/1.20.3/libexec/src/runtime/debug/stack.go:24 +0x64
...
}
```
While kytsya fetching panic, panic message and a stack trace returns as a normal error!

3. Needs to handle every list member in a separate goroutine?
```
	data := []int{1, 2, 3}
	resCh := NewEachRunner[int, string](data).
		Handle(func(val int) Result[string] {
			return Result[string]{Data: fmt.Sprint(val)}
		})

	ForChan(resCh, func(val Result[string]) {
		fmt.Println(val.Data, val.Err)
	})
```
Every member of "data" slice will be handled in a separate goroutine ("handler") and results will be returned in resCh channel that will be closed after all tasks are done.

No external dependencies in the kit, pure std golang :)
