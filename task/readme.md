### what is Task:

`Task` is a wrap for `sync.WatiGroup`, its aim is to make `go func` code more beautiful and easy to read

### notice:

you ***`must`*** initialize your variable including all fields in variable before run tasks

### how to use:

there are two ways to use task:

- RunTasks, use `Runner` in Task:
  - like "out" keyword in c#, and you must pass value to it yourself
  - example: [runTasks](https://github.com/mint-leaf/gotokit/blob/master/task/task_test.go#L32)
- RunFunc, use `Func` in Task:
  - uses reflect to pass value, it runs more slowly, but it is more simple
  - example: [runFunc](https://github.com/mint-leaf/gotokit/blob/master/task/task_test.go#L54)

### performance

cpu: AMD Ryzen 7 5800X 8-Core Processor
+ BenchmarkRunFunc
    + BenchmarkRunFunc-16               616291              1724 ns/op
+ BenchmarkRunTasks
    + BenchmarkRunTasks-16              747738              1491 ns/op
+ BenchmarkRunErrGroup
    + BenchmarkRunErrGroup-16           799855              1382 ns/op
