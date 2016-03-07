package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

	plug := PluginManager{
		configFile: filepath.Join(p, "test_config.json"),
		name:       name,
		path:       filepath.Join(p, name),
	}
	plug.Start()
	defer plug.Stop()
	err = plug.Configure()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Finish configure")

	err = plug.Main()
	if err != nil {
		fmt.Println(err)
		return
	}
	plug.Notify()
	fmt.Println("Finish")
}
