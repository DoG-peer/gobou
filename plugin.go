package main

import (
	"github.com/DoG-peer/gobou/utils"
	"github.com/dullgiulio/pingo"
	"path/filepath"
	"runtime"
	"time"
)

// PluginManager stores the infomation of the running plugin
type PluginManager struct {
	interval       time.Duration
	err            error
	configFile     string
	name           string
	path           string
	dataDir        string
	plugin         *pingo.Plugin
	pluginInfo     PluginInfo
	messages       []gobou.Message
	messageChannel chan<- gobou.Message
}

// PluginInfo is saved in the main config.json
type PluginInfo struct {
	Name            string
	Path            string
	Repository      string
	SourceDirectory string
	CacheDirectory  string
	DataDirectory   string
	ConfigFile      string
}

// MakePluginInfo called by app
func (app *AppPath) MakePluginInfo(repository string, name string) PluginInfo {
	bin := name
	if runtime.GOOS == "windows" {
		bin += ".bin"
	}
	return PluginInfo{
		Name:            name,
		Path:            filepath.Join(app.PluginDir, bin),
		Repository:      repository,
		SourceDirectory: filepath.Join(app.CacheDir, "src", name),
		CacheDirectory:  filepath.Join(app.CacheDir, "plugin", name),
		DataDirectory:   filepath.Join(app.DataDir, "plugin", name),
		ConfigFile:      filepath.Join(app.ConfigDir, "plugin_config", name+".json"),
	}
}

// Load plugin by manager
func (p *PluginManager) Load(plug PluginInfo, mesChan chan<- gobou.Message) {
	p.pluginInfo = plug
	p.name = plug.Name
	p.path = plug.Path
	p.configFile = plug.ConfigFile
	p.dataDir = plug.DataDirectory
	p.messageChannel = mesChan
}

// Start starts Task
func (p *PluginManager) Start() {
	if runtime.GOOS == "windows" {
		p.plugin = pingo.NewPlugin("tcp", p.path)
	} else {
		p.plugin = pingo.NewPlugin("unix", p.path)
	}
	p.plugin.Start()
}

// Stop stops Task
func (p *PluginManager) Stop() {
	p.plugin.Stop()
}

// Main calls main plug
func (p *PluginManager) Main() error {
	return p.plugin.Call("Task.Main", "", &p.messages)
}

// Notify notifies the information from plugin to desktop
func (p *PluginManager) Notify() error {
	for _, mes := range p.messages {
		if mes.IsNone() {
			continue
		}
		p.messageChannel <- mes
	}
	return nil
}

// SaveData saves data
func (p *PluginManager) SaveData() error {
	return p.plugin.Call("Task.SaveData", p.dataDir, &p.err)
}

// SaveConfig saves config
func (p *PluginManager) SaveConfig() error {
	return p.plugin.Call("Task.SaveConfig", p.configFile, &p.err)
}

// ReadInterval checks wait interval
func (p *PluginManager) ReadInterval() error {
	return p.plugin.Call("Task.Interval", "", &p.interval)
}

// Configure load config
func (p *PluginManager) Configure() error {
	return p.plugin.Call("Task.Configure", p.configFile, &p.err)
}

// Wait waits until next step
func (p *PluginManager) Wait() error {
	if err := p.ReadInterval(); err != nil {
		return err
	}
	time.Sleep(p.interval)
	return nil
}
