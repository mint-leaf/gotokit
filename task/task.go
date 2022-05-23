package task

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
)

// Task is represent as the basic element of schedule
type Task struct {
	Receiver interface{}
	Runner   func(out interface{}) error
	Func     func() (interface{}, error)
}

// EmptyReceiver is used for func like `func(a int) error` which has only error as return value or has no return value
// usually, multi return value func use `Func`, single or zero return value func use `Runner`
// but if after func return you want to pass value, maybe Runner is better
var (
	EmptyReceiver = uint32(0)
	defaultErrOut = os.Stderr
	errOut        io.Writer
)

func SetErrorOutput(w io.Writer) {
	errOut = w
}

func ErrorOutput() io.Writer {
	if errOut == nil {
		return defaultErrOut
	}
	return errOut
}

// Check is used to check if t.receiver is ptr
func (t Task) Check() error {
	ty := reflect.TypeOf(t.Receiver)
	if ty.Kind() != reflect.Ptr {
		return fmt.Errorf("invalid type: %s", ty.Name())
	}
	return nil
}

// Run is used for Func,
func (t Task) Run() error {
	if err := t.Check(); err != nil {
		return err
	}
	v, err := t.Func()
	if err != nil {
		return err
	}
	reflect.ValueOf(t.Receiver).Elem().Set(reflect.ValueOf(v).Elem())
	return nil
}

func RunTasks(tasks ...Task) error {
	if len(tasks) == 0 {
		return nil
	}
	if len(tasks) == 1 {
		if err := tasks[0].Check(); err != nil {
			return err
		}
		return tasks[0].Runner(tasks[0].Receiver)
	}
	wg, hasErr := new(sync.WaitGroup), uint32(0)
	wg.Add(len(tasks))
	for index := range tasks {
		go func(t *Task, index uint32) {
			defer func() {
				if err := recover(); err != nil {
					_, _ = fmt.Fprintln(ErrorOutput(), err)
				}
				wg.Done()
			}()
			if err := t.Check(); err != nil {
				atomic.CompareAndSwapUint32(&hasErr, 0, index)
			} else if err = t.Runner(t.Receiver); err != nil {
				atomic.CompareAndSwapUint32(&hasErr, 0, index)
			}
		}(&tasks[index], uint32(index))
	}
	wg.Wait()
	if hasErr != 0 {
		return fmt.Errorf("error first occurs index: %d", hasErr)
	}
	return nil
}

func RunFunc(tasks ...Task) error {
	if len(tasks) == 0 {
		return nil
	}
	if len(tasks) == 1 {
		if err := tasks[0].Check(); err != nil {
			return err
		}
		return tasks[0].Run()
	}
	wg, hasErr := new(sync.WaitGroup), uint32(0)
	wg.Add(len(tasks))
	for index := range tasks {
		go func(t *Task, index uint32) {
			defer func() {
				if err := recover(); err != nil {
					_, _ = fmt.Fprintln(ErrorOutput(), err)
				}
				wg.Done()
			}()
			if err := t.Run(); err != nil {
				atomic.CompareAndSwapUint32(&hasErr, 0, index)
			}
		}(&tasks[index], uint32(index))
	}
	wg.Wait()
	if hasErr != 0 {
		return fmt.Errorf("error first occurs index: %d", hasErr)
	}
	return nil
}
