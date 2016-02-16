package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

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

// AppConfig stores config about this application
type AppConfig struct {
	Plugins []PluginInfo
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
	home := GetHome()

	xdgConfigDir := os.Getenv("XDG_CONFIG_DIR")
	xdgDataDir := os.Getenv("XDG_DATA_DIR")
	xdgCacheHome := os.Getenv("XDG_CACHE_HOME")

	if xdgConfigDir == "" {
		xdgConfigDir = filepath.Join(home, ".config")
	}
	if xdgDataDir == "" {
		xdgDataDir = filepath.Join(home, ".local/share")
	}
	if xdgCacheHome == "" {
		xdgCacheHome = filepath.Join(home, ".cache")
	}

	configDir := filepath.Join(xdgConfigDir, name)
	return AppPath{
		ConfigDir:       configDir,
		DataDir:         filepath.Join(xdgDataDir, name),
		CacheDir:        filepath.Join(xdgCacheHome, name),
		ConfigFile:      filepath.Join(configDir, "config.json"),
		PluginDir:       filepath.Join(configDir, "plugin"),
		PluginConfigDir: filepath.Join(configDir, "plugin_config"),
	}

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

// Add plugin to config
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

// SaveConfig saves config
func (app *AppPath) SaveConfig() {
	ioutil.WriteFile(app.ConfigFile, []byte(app.Config.String()), os.ModePerm)
}
