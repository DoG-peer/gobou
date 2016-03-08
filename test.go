package main

import (
	"fmt"
	"github.com/DoG-peer/gobou/notify"
	"github.com/DoG-peer/gobou/utils"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Test runs test for gobou plugin developer
func Test() {
	p, _ := os.Getwd()
	name := filepath.Base(p)
	fmt.Println(p)

	err := exec.Command("go", "build").Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Finish build")

	mesChan := make(chan gobou.Message)
	go func() {
		time.Sleep(10 * time.Second)
		close(mesChan)
	}()

	plug := PluginManager{
		configFile:     filepath.Join(p, "test_config.json"),
		name:           name,
		path:           filepath.Join(p, name),
		messageChannel: mesChan,
	}
	plug.Start()
	defer plug.Stop()
	err = plug.Configure()
	if err != nil {
		fmt.Println(err, "config")
		return
	}
	fmt.Println("Finish configure")

	go func() {
		err = plug.Main()
		if err != nil {
			fmt.Println(err, "plugin")
			return
		}
		plug.Notify()
	}()

	for mes := range mesChan {
		if !mes.NotifyMessage.IsNone() {
			notify.Notify(mes.NotifyMessage.Text)
		}
		if !mes.StdoutMessage.IsNone() {
			fmt.Println(mes.StdoutMessage.Text)
		}
	}

	fmt.Println("Finish")
}
