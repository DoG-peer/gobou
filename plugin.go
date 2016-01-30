package main

import (
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
}
