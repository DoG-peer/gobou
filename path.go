package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// AppConfig stores config about this application
type AppConfig struct {
	Plugins []PluginInfo
}

// AppPath stores paths about this application
type AppPath struct {
	ConfigDir       string
	DataDir         string
	CacheDir        string
	PluginDir       string
	PluginConfigDir string
	ConfigFile      string
	Config          AppConfig
}

// PrepareDirs make directories and so on
func (app *AppPath) PrepareDirs() error {
	dirs := []string{
		app.ConfigDir,
		app.DataDir,
		app.CacheDir,
		app.PluginDir,
		app.PluginConfigDir,
	}

	for _, dir := range dirs {
		if finfo, e := os.Stat(dir); os.IsNotExist(e) {
			err := os.MkdirAll(dir, 0777)
			if err != nil {
				return err
			}
		} else if !finfo.IsDir() {
			return fmt.Errorf("%s is not directory", dir)
		}

	}

	// app.configFile
	if finfo, e := os.Stat(app.ConfigFile); !os.IsNotExist(e) && finfo.IsDir() {
		return fmt.Errorf("%s is directory", app.ConfigFile)
	}
	return nil

}

// GetHome returns absolute path of home directory
func GetHome() string {
	home := os.Getenv("HOME")
	if home == "" {
		home = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"))
	}
	return home
}

// GetAppPath initialize application
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

	app.ConfigFile = filepath.Join(app.ConfigDir, "config.json")
	app.PluginDir = filepath.Join(app.ConfigDir, "plugin")
	app.PluginConfigDir = filepath.Join(app.ConfigDir, "plugin_config")

	return app
}

// Configure loads config file
func (app *AppPath) Configure() error {
	file, err := ioutil.ReadFile(app.ConfigFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	err = json.Unmarshal(file, &app.Config)
	if err != nil {
		return err
	}
	return nil
}

func (c *AppConfig) Add(plug PluginInfo) {
	newConfig := []PluginInfo{plug}
	for _, p := range c.Plugins {
		if p.Name != plug.Name {
			newConfig = append(newConfig, p)
		}
	}
	c.Plugins = newConfig
}

func (c *AppConfig) String() string {
	s, _ := json.MarshalIndent(c, "", "  ")
	return string(s)
}

func (app *AppPath) SaveConfig() {
	ioutil.WriteFile(app.ConfigFile, []byte(app.Config.String()), os.ModePerm)
}
