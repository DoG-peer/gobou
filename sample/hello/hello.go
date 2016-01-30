package main

import (
	"github.com/dullgiulio/pingo"
	"os/exec"
	"strconv"
	"time"
)

type AppTask interface {
	run()
	configure()
	interval() time.Duration
	self() *AppTask
}

// Edit your task
type Task struct {
	count int
}

func (p *Task) Configure(configFile string, e *error) error {
	return nil
}
func (p *Task) Main(configFile string, e *error) error {
	p.count++
	_, *e = exec.Command("notify-send", "test"+strconv.Itoa(p.count)).Output()
	return *e
}
func (p *Task) SaveData(configFile string, e *error) error {
	return nil
}
func (p *Task) SaveConfig(configFile string, e *error) error {
	return nil
}
func (p *Task) Interval(a string, d *time.Duration) error {
	*d = 1 * time.Second
	return nil
}
func (p *Task) End() error {
	return nil
}
func main() {
	task := &Task{count: 0}
	pingo.Register(task)
	pingo.Run()
}
