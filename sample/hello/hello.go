package main

import (
	"github.com/DoG-peer/gobou/utils"
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
func (p *Task) Main(configFile string, s *[]gobou.Message) error {
	p.count++
	*s = []gobou.Message{
		gobou.Notify(strconv.Itoa(p.count) + "\nhello by notify"),
		gobou.Print(strconv.Itoa(p.count) + "\nhello by print"),
		gobou.Say(strconv.Itoa(p.count)),
	}
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
	gobou.Register(&Task{count: 0})
}
