package main

import (
	"errors"
	"fmt"
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
	configDir       string
	dataDir         string
	cacheDir        string
	pluginDir       string
	pluginConfigDir string
	configFile      string
	config          *AppConfig
}

func (app *AppPath) prepareDirs() error {
	// app.configDir
	if finfo, e := os.Stat(app.configDir); os.IsNotExist(e) {
		err := os.Mkdir(app.configDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.configDir)
	}

	// app.dataDir
	if finfo, e := os.Stat(app.dataDir); os.IsNotExist(e) {
		err := os.Mkdir(app.dataDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.dataDir)
	}

	// app.cacheDir
	if finfo, e := os.Stat(app.cacheDir); os.IsNotExist(e) {
		err := os.Mkdir(app.cacheDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.cacheDir)
	}

	// app.pluginDir
	if finfo, e := os.Stat(app.pluginDir); os.IsNotExist(e) {
		err := os.Mkdir(app.pluginDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.pluginDir)
	}

	// app.pluginConfigDir
	if finfo, e := os.Stat(app.pluginConfigDir); os.IsNotExist(e) {
		err := os.Mkdir(app.pluginConfigDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.pluginConfigDir)
	}

	// app.configFile
	if finfo, e := os.Stat(app.configFile); !os.IsNotExist(e) && finfo.IsDir() {
		return fmt.Errorf("%s is directory", app.configFile)
	}
	return nil

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
		plugins = append(plugins, filepath.Join(app.pluginDir, plugin.Name()))
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
	app.pluginConfigDir = filepath.Join(app.configDir, "plugin_config")

	return app
}

//TODO
func (app *AppPath) configure() error {
	return nil
}

/*
func (app *AppPath) configurePlugin(task *AppTask, pluginPath string) error {
	return nil
}
*/

// Task stores the infomation of the running plugin
type Task struct {
	interval   time.Duration
	err        error
	configFile string
	name       string
	path       string
	dataFile   string
}

func main() {
	app := getAppPath("gobou")
	if err := app.prepareDirs(); err != nil {
		log.Fatal(err)
		return
	}
	pluginPaths, perr := app.getPlugins()
	if perr != nil {
		log.Fatal(perr)
		return
	}
	plugins := []*pingo.Plugin{}
	//tasks := []*AppTask{}
	for _, pluginPath := range pluginPaths {
		log.Println(pluginPath)

		plug := pingo.NewPlugin("unix", pluginPath)
		plugins = append(plugins, plug)
		go func() {
			task := Task{}
			task.path = pluginPath
			task.name = filepath.Base(pluginPath)
			task.configFile = filepath.Join(app.pluginConfigDir, task.name)
			task.dataFile = filepath.Join(app.dataDir, task.name)

			plug.Start()
			defer plug.Stop()
			if err := plug.Call("Task.Configure", task.configFile, &task.err); err != nil {
				log.Fatalf("Failed to configurePlugin %s", task.configFile)
				return
			}
			for {
				// main task
				if err := plug.Call("Task.Main", "", &task.err); err != nil {
					log.Fatalln(err)
					break
				}
				// log
				if err := plug.Call("Task.SaveData", task.dataFile, &task.err); err != nil {
					log.Fatalln(err)
					break
				}
				// change config
				if err := plug.Call("Task.SaveConfig", task.configFile, &task.err); err != nil {
					log.Fatalln(err)
					break
				}

				// wait
				if err := plug.Call("Task.Interval", "", &task.interval); err != nil {
					log.Fatalln(err)
					break
				} else {
					time.Sleep(task.interval)
				}
			}
		}()
	}
	time.Sleep(10 * time.Second)

}
