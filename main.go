package main

import (
	"errors"
	"github.com/dullgiulio/pingo"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

// AppConfig stores config about this application
type AppConfig struct {
}

// AppPath stores paths about this application
type AppPath struct {
	configDir, dataDir, cacheDir, pluginDir, configFile string
	config                                              *AppConfig
}

//TODO
func (app *AppPath) prepareDirs() {
}

func (app *AppPath) getPlugins() ([]string, error) {
	plugins := []string{}
	if app.pluginDir == "" {
		return nil, errors.New("your AppPath does not have a plugin path")
	}
	finfo, e := os.Stat(app.pluginDir)
	if e != nil {
		return nil, e
	}
	if !finfo.IsDir() {
		return nil, errors.New(app.pluginDir + " is not a directory")
	}
	files, e2 := ioutil.ReadDir(app.pluginDir)

	if e2 != nil {
		return nil, e2
	}

	for _, plugin := range files {
		plugins = append(plugins, plugin.Name())
	}
	return plugins, nil
}

func getAppPath(name string) AppPath {
	app := AppPath{}

	home := os.Getenv("HOME")
	var xdgConfigDir, xdgDataDir, xdgCacheHome string
	if os.Getenv("XDG_CONFIG_DIR") != "" {
		xdgConfigDir = os.Getenv("XDG_CONFIG_DIR")
	} else {
		xdgConfigDir = filepath.Join(home, ".config")
	}

	app.configDir = filepath.Join(xdgConfigDir, name)

	if os.Getenv("XDG_DATA_DIR") != "" {
		xdgDataDir = os.Getenv("XDG_CONFIG_DIR")
	} else {
		xdgDataDir = filepath.Join(home, ".local/share")
	}

	app.dataDir = filepath.Join(xdgDataDir, name)

	if os.Getenv("XDG_CACHE_HOME") != "" {
		xdgCacheHome = os.Getenv("XDG_CACHE_HOME")
	} else {
		xdgCacheHome = filepath.Join(home, ".cache")
	}

	app.cacheDir = filepath.Join(xdgCacheHome, name)

	app.configFile = filepath.Join(app.configDir, "config.yml")
	app.pluginDir = filepath.Join(app.configDir, "plugin")

	return app
}

//TODO
func (app *AppPath) configure() error {
	return nil
}
func (app *AppPath) configurePlugin(task *AppTask, pluginPath string) error {
	return nil
}

// AppTask stores a task defined in a plugin
type AppTask interface {
	run()
	configure()
	interval() time.Duration
}

func main() {
	app := getAppPath("gobou")
	app.prepareDirs()
	pluginPaths, perr := app.getPlugins()
	if perr != nil {
		log.Fatal(perr)
		return
	}
	plugins := []*pingo.Plugin{}
	tasks := []*AppTask{}
	for _, pluginPath := range pluginPaths {
		plug := pingo.NewPlugin("unix", pluginPath)
		plugins = append(plugins, plug)

		plug.Start()
		var task *AppTask
		err := plug.Call("Task", nil, task)
		if err != nil {
			log.Fatalln(err, "On loading "+pluginPath)
			continue
		}
		errConf := app.configurePlugin(task, pluginPath)
		if errConf != nil {
			log.Fatalln(err, "On loading "+pluginPath)
			continue
		}
		tasks = append(tasks, task)
		plug.Stop()
	}

	for _, task := range tasks {
		go func() {
			for {
				(*task).run()
				time.Sleep((*task).interval())
			}
		}()
	}
	log.Println(app)
	for {
	}
}
