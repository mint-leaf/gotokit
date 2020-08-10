package task

import (
	"fmt"
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
			if err := t.Check(); err != nil {
				atomic.CompareAndSwapUint32(&hasErr, 0, index)
			} else if err = t.Runner(t.Receiver); err != nil {
				atomic.CompareAndSwapUint32(&hasErr, 0, index)
			}
			wg.Done()
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
			if err := t.Run(); err != nil {
				atomic.CompareAndSwapUint32(&hasErr, 0, index)
			}
			wg.Done()
		}(&tasks[index], uint32(index))
	}
	wg.Wait()
	if hasErr != 0 {
		return fmt.Errorf("error first occurs index: %d", hasErr)
	}
	return nil
}
