package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Configurator is used by gobou config command
type Configurator struct {
	isMain bool
	plugin string
}

func edit(fname string) {
	_, e := os.Stat(fname)
	if os.IsNotExist(e) {
		log.Fatalln(e)
		return
	}
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", fname).Run()
	case "windows":
		exec.Command("notepad", fname).Run()
	case "mac":
		exec.Command("open", fname).Run()
	}
}

// OpenMainConfig opens config file of gobou
func (c *Configurator) OpenMainConfig(confFile string) {
	_, e := os.Stat(confFile)
	if os.IsNotExist(e) {
		ioutil.WriteFile(confFile, []byte("{}"), os.ModePerm)
	}
	edit(confFile)
}

// OpenPluginConfig opens config file of the plugin
func (c *Configurator) OpenPluginConfig(pluginConfigDir string) {
	file := filepath.Join(pluginConfigDir, c.plugin+".json")
	edit(file)
}
