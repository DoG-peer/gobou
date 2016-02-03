package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// AppConfig stores config about this application
type AppConfig struct {
}

// AppPath stores paths about this application
type AppPath struct {
	ConfigDir       string
	DataDir         string
	CacheDir        string
	PluginDir       string
	PluginConfigDir string
	ConfigFile      string
	Config          *AppConfig
}

func (app *AppPath) PrepareDirs() error {
	// app.configDir
	if finfo, e := os.Stat(app.ConfigDir); os.IsNotExist(e) {
		err := os.MkdirAll(app.ConfigDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.ConfigDir)
	}

	// app.dataDir
	if finfo, e := os.Stat(app.DataDir); os.IsNotExist(e) {
		err := os.MkdirAll(app.DataDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.DataDir)
	}

	// app.cacheDir
	if finfo, e := os.Stat(app.CacheDir); os.IsNotExist(e) {
		err := os.MkdirAll(app.CacheDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.CacheDir)
	}

	// app.pluginDir
	if finfo, e := os.Stat(app.PluginDir); os.IsNotExist(e) {
		err := os.MkdirAll(app.PluginDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.PluginDir)
	}

	// app.pluginConfigDir
	if finfo, e := os.Stat(app.PluginConfigDir); os.IsNotExist(e) {
		err := os.MkdirAll(app.PluginConfigDir, 0777)
		if err != nil {
			return err
		}
	} else if !finfo.IsDir() {
		return fmt.Errorf("%s is not directory", app.PluginConfigDir)
	}

	// app.configFile
	if finfo, e := os.Stat(app.ConfigFile); !os.IsNotExist(e) && finfo.IsDir() {
		return fmt.Errorf("%s is directory", app.ConfigFile)
	}
	return nil

}

func (app *AppPath) GetPlugins() ([]string, error) {
	plugins := []string{}
	if app.PluginDir == "" {
		return nil, errors.New("your AppPath does not have a plugin path")
	}
	finfo, e := os.Stat(app.PluginDir)
	if e != nil {
		return nil, e
	}
	if !finfo.IsDir() {
		return nil, errors.New(app.PluginDir + " is not a directory")
	}
	files, e2 := ioutil.ReadDir(app.PluginDir)

	if e2 != nil {
		return nil, e2
	}

	for _, plugin := range files {
		plugins = append(plugins, filepath.Join(app.PluginDir, plugin.Name()))
	}
	return plugins, nil
}

func GetHome() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
	}
	return home
}

func GetAppPath(name string) AppPath {
	app := AppPath{}

	home := GetHome()
	var xdgConfigDir, xdgDataDir, xdgCacheHome string
	if os.Getenv("XDG_CONFIG_DIR") != "" {
		xdgConfigDir = os.Getenv("XDG_CONFIG_DIR")
	} else {
		xdgConfigDir = filepath.Join(home, ".config")
	}

	app.ConfigDir = filepath.Join(xdgConfigDir, name)

	if os.Getenv("XDG_DATA_DIR") != "" {
		xdgDataDir = os.Getenv("XDG_CONFIG_DIR")
	} else {
		xdgDataDir = filepath.Join(home, ".local/share")
	}

	app.DataDir = filepath.Join(xdgDataDir, name)

	if os.Getenv("XDG_CACHE_HOME") != "" {
		xdgCacheHome = os.Getenv("XDG_CACHE_HOME")
	} else {
		xdgCacheHome = filepath.Join(home, ".cache")
	}

	app.CacheDir = filepath.Join(xdgCacheHome, name)

	app.ConfigFile = filepath.Join(app.ConfigDir, "config.yml")
	app.PluginDir = filepath.Join(app.ConfigDir, "plugin")
	app.PluginConfigDir = filepath.Join(app.ConfigDir, "plugin_config")

	return app
}

//TODO
func (app *AppPath) Configure() error {
	return nil
}

/*
func (app *AppPath) configurePlugin(task *AppTask, pluginPath string) error {
	return nil
}
*/
