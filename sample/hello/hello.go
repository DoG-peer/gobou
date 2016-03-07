package main

import (
	"github.com/dullgiulio/pingo"
	"os/exec"
	"strconv"
	"time"
)

// Task manages the activity of your plugin
// Task's methods must be func[string, *inteface{}] error
// you can not change this name
type Task struct {
	count int
}

func (p *Task) Configure(configFile string, e *error) error {
	return nil
}
func (p *Task) Main(configFile string, s *string) error {
	p.count++
	*s = strconv.Itoa(p.count) + "\nhello"
	return nil
}
func (p *Task) SaveData(configFile string, e *error) error {
	return nil
}
func (p *Task) SaveConfig(configFile string, e *error) error {
	return nil
}
func (p *Task) Interval(a string, d *time.Duration) error {
	*d = 10 * time.Second
	return nil
}
func (p *Task) End(a string, b *interface{}) error {
	exec.Command("notify-send", "hoge").Run()
	return nil
}
func main() {
	task := &Task{count: 0}
	pingo.Register(task)
	pingo.Run()
}
