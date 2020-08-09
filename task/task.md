### what is Task:

`Task` is a wrap for `sync.WatiGroup`, its aim is to make `go func` code more beautiful and easy to read

### notice:

you $\color{#FF0000}`must`$ initialize your variable including all fields in variable before run tasks

### how to use:

there are two ways to use task:

- RunTasks, use `Runner` in Task:
  - like "out" keyword in c#, and you must pass value to it yourself
  - example: [runTasks](./task_test.go)
- RunFunc, use `Func` in Task:
  - uses reflect to pass value, it runs more slowly, but it is more simple
  - example: [runFunc](./task_test.go)

### performance

- RunFunc
  - | BenchmarkRunFunc | 353136 | 3548 ns/op
- RunTasks
  - | BenchmarkRunTasks | 500010 | 2652 ns/op
