package main

import (
	"github.com/dullgiulio/pingo"
	"path/filepath"
	"time"
)

// Task stores the infomation of the running plugin
type Task struct {
	interval   time.Duration
	err        error
	configFile string
	name       string
	path       string
	dataFile   string
	plugin     *pingo.Plugin
}

// LoadTask initialize Task
func LoadTask(pluginPath, pluginConfigDir, dataDir string) Task {
	task := Task{}
	task.path = pluginPath
	task.name = filepath.Base(pluginPath)
	task.configFile = filepath.Join(pluginConfigDir, task.name)
	task.dataFile = filepath.Join(dataDir, task.name)
	return task
}

// Start starts Task
func (t *Task) Start() {
	t.plugin = pingo.NewPlugin("unix", t.path)
	t.plugin.Start()
}

// Stop stops Task
func (t *Task) Stop() {
	t.plugin.Stop()
}

// Main calls main task
func (t *Task) Main() error {
	return t.plugin.Call("Task.Main", "", &t.err)
}

// SaveData saves data
func (t *Task) SaveData() error {
	return t.plugin.Call("Task.SaveData", t.dataFile, &t.err)
}

// SaveConfig saves config
func (t *Task) SaveConfig() error {
	return t.plugin.Call("Task.SaveConfig", t.configFile, &t.err)
}

// ReadInterval checks wait interval
func (t *Task) ReadInterval() error {
	return t.plugin.Call("Task.Interval", "", &t.interval)
}

// Configure load config
func (t *Task) Configure() error {
	return t.plugin.Call("Task.Configure", t.configFile, &t.err)
}

// Wait waits until next step
func (t *Task) Wait() error {
	if err := t.ReadInterval(); err != nil {
		return err
	}
	time.Sleep(t.interval)
	return nil
}
