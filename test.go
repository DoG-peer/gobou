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
	err := exec.Command("go", "build").Run()
	if err != nil {
		fmt.Println(err)
		return
	}
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
	err = plug.Main()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p)
}
