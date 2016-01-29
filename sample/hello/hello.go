package main

import (
	"github.com/dullgiulio/pingo"
	"strconv"
	"time"
)

// Edit your task
type Task struct {
	count int
}

func (p *Task) Foo(a string, msg *string) error {
	p.count++
	//fmt.Printf("Hello world#%d\n", p.count)
	*msg = "Hello world" + strconv.Itoa(p.count)
	return nil
}

func (p *Task) Interval(a string, d *time.Duration) error {
	*d = 1 * time.Second
	return nil
}

func (p *Task) Self() *Task {
	return p
}

func main() {
	task := &Task{count: 0}
	pingo.Register(task)
	pingo.Run()
}
