# 🇺🇦kytsya🇺🇦
![Logo](https://github.com/bkatrenko/kytsya/blob/master/logo.png)

## Go toolkit!
Kytsya means "kitten" in ukrainian. While cats are the best programmers friends, we decided to call the repo in honour of these beautiful creatures.
It's a small, but powerful kit that could be used to change an approach to threading and working with slices.

It contains:
* Number of controllers for goroutines to use them in a functional, save and readable manner
* Controllers for run goroutines that should return a result
* Builders for start a goroutines with a wait group, recovery handler or/and timeout
* ForEach, ForEachChan, Map, Reduce, Filter functions to handle data in a functional and readable manner

Cases where kytsya will shine:
* Work hard with goroutines? Kytsya helps to manage that in nice way.
* Need recovery handlers, but don't want to manually control every case that should be covered with deferred recoveries? Here you could find __WithRecover()__ builder, that makes recovery graceful!
* Run goroutins that returns results / need to make everything stable and built in one manner? Kytsya will help to do it with __ErrorBox__ or __ForEachRunner__.
* Every developer in a big team uses its own threading style? Have a mix of channels, waitgroups, "res" slices? Kytsya gives one approach for everything.
* Work with data? __ForEach__, __Map__, __Reduce__ will definetly what is missing in go's std lib.

Kytsya is not a framework, it's a toolkit that gives an ability, but not forces developers to use the one way.
Here is only small part from what it is doing:
1. Need to run a set of goroutines with recovery handlers and a wait group?
```
    kytsya.NewBox().
   	 WithRecover().
   	 WithWaitGroup().
   	 AddTask(func() { fmt.Println("🐈") }).
   	 AddTask(func() { fmt.Println("🐈") }).
   	 AddTask(func() { fmt.Println("🐈") }).
   	 Run().Wait()
```
2. Need a safe way to run a set of goroutines and properly read some results?
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

3. Need to handle every list member in a separate goroutine?
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
Every member of the "data" slice will be handled in a separate goroutine ("handler") and results will be returned in the resCh channel that will be closed after all tasks are done.

4. Need Map, Reduce, ForEach and Filter functions?
```
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
```
Here it is!

### See example_test.go for more examples!
### We appreciate feedbacks and found issues! Feel free to become a contributor and add here things you miss in go as we do!

## Why should we use kytsya?
- It is giving a way to make an application that work hard with a data to handle everything in more functional way
- One toolkit for all: big projects struggling from different threading styles of each team member. With kytsya it is not a problem anymore!
- It is reliable and well-tested
- Kytsya gives an ability to handle things graceful, but not force - it is not necessary to use whole toolkit
- The way to do a beautiful things quick and graceful
- Most of its controller implements laziness - no work until result is requested (but not them all :))

## Benchmarking
Benchmarks included in the repo.
Here is the results for the ErrTaskRunner (handler for a group of goroutines that returns a results):
```
goos: darwin
goarch: arm64
pkg: github.com/bkatrenko/kytsya
BenchmarkErrorBox/pure_Go-10         	  524178	      2210 ns/op	     360 B/op	       8 allocs/op
BenchmarkErrorBox/kytsunya-10        	  443697	      2333 ns/op	     512 B/op	      12 allocs/op
PASS
ok  	github.com/bkatrenko/kytsya	3.584s
```

## also:
- No external dependencies in the kit, pure std golang :)
- 100% test coverage
