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
		time.Sleep(5 * time.Second)
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
		fmt.Println(err)
		return
	}
	fmt.Println("Finish configure")

	go func() {
		err = plug.Main()
		if err != nil {
			fmt.Println(err)
			return
		}
		plug.Notify()
	}()

	for mes := range mesChan {
		if mes.IsNone() {
			continue
		}
		err := notify.Notify(mes.NotifyMessage.Text)
		if err != nil {
			fmt.Println(err)
			break
		}
	}

	fmt.Println("Finish")
}
