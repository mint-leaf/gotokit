package task

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"golang.org/x/sync/errgroup"
)

type a struct {
	I int
}

type b struct {
	S string
}

type c struct {
	V *a
}

type d struct {
	A1 *a
	B1 *b
	C1 *c
	E1 []int
}

var d1 = d{A1: new(a), B1: new(b), C1: new(c), E1: make([]int, 0)}

func runTasks() error {
	return RunTasks(
		Task{Receiver: d1.A1, Runner: func(out interface{}) error {
			*out.(*a) = a{I: 1}
			return nil
		}},
		Task{Receiver: d1.B1, Runner: func(out interface{}) error {
			*out.(*b) = b{S: "b"}
			return nil
		}},
		Task{Receiver: d1.C1, Runner: func(out interface{}) error {
			*out.(*c) = c{V: &a{I: 2}}
			return nil
		}},
		Task{Receiver: &d1.E1, Runner: func(out interface{}) error {
			*out.(*[]int) = []int{1}
			return nil
		}},
	)
}

func runFunc() error {
	return RunFunc(
		Task{Receiver: d1.A1, Func: func() (out interface{}, err error) {
			return &a{I: 1}, nil
		}},
		Task{Receiver: d1.B1, Func: func() (out interface{}, err error) {
			return &b{S: "b"}, nil
		}},
		Task{Receiver: d1.C1, Func: func() (out interface{}, err error) {
			return &c{V: &a{I: 2}}, nil
		}},
		Task{Receiver: &d1.E1, Func: func() (out interface{}, err error) {
			return &[]int{1}, nil
		}},
	)
}

func runErrGroup() error {
	g := new(errgroup.Group)
	g.Go(func() error {
		d1.A1 = &a{I: 1}
		return nil
	})
	g.Go(func() error {
		d1.B1 = &b{S: "b"}
		return nil
	})
	g.Go(func() error {
		d1.C1 = &c{V: &a{I: 2}}
		return nil
	})
	g.Go(func() error {
		d1.E1 = []int{1}
		return nil
	})
	return g.Wait()
}

func check() {
	if d1.A1.I != 1 {
		log.Fatalln("a value invalid")
	}
	if d1.B1.S != "b" {
		log.Fatalln("b value invalid")
	}
	if d1.C1.V.I != 2 {
		log.Fatalln("c value invalid")
	}
	if len(d1.E1) != 1 && d1.E1[0] != 1 {
		log.Fatalln("e value invalid")
	}
}

func TestRunTasks(t *testing.T) {
	err := runTasks()
	if err != nil {
		log.Fatal(err)
	}
	check()
	dv, _ := json.Marshal(d1)
	fmt.Println(string(dv))
}

func TestRunFunc(t *testing.T) {
	err := runFunc()
	if err != nil {
		log.Fatal(err)
	}
	check()
	dv, _ := json.Marshal(d1)
	fmt.Println(string(dv))
}

func BenchmarkRunFunc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = runFunc()
	}
}

func BenchmarkRunTasks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = runTasks()
	}
}

func BenchmarkRunErrGroup(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = runErrGroup()
	}
}
