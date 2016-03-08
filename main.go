package main

import (
	"fmt"
	"github.com/DoG-peer/gobou/utils"
	"log"
	"time"
)

func main() {
	// check arguments
	cli := ParseCliInfo()
	// store pathdata
	app := GetAppPath("gobou")
	// load configfile
	if err := app.Configure(); err != nil {
		log.Fatal(err)
		return
	}
	// make directories
	if err := app.PrepareDirs(); err != nil {
		log.Fatal(err)
		return
	}

	// switch to each command
	switch {
	case cli.isHelp:
		cli.ShowHelp()
		return
	case cli.isInstall:
		cli.installInfo.Install(app.PluginDir, app.CacheDir, app.PluginConfigDir)
		plug := app.MakePluginInfo(cli.installInfo.url, cli.installInfo.name)
		app.Config.Add(plug)
		app.SaveConfig()
		return
	case cli.isUpdate:
		cli.installInfo.Update(app.PluginDir, app.CacheDir, app.PluginConfigDir)
		plug := app.MakePluginInfo(cli.installInfo.url, cli.installInfo.name)
		app.Config.Add(plug)
		app.SaveConfig()
		return
	case cli.isTest:
		Test()
		return
	case cli.isGenerate:
		log.Println("generate")
		return
	case cli.isConfig:
		if cli.configInfo.isMain {
			cli.configInfo.OpenMainConfig(app.ConfigFile)
		} else {
			cli.configInfo.OpenPluginConfig(app.PluginConfigDir)

		}
		return
	}

	// search plugins
	plugins := app.Config.Plugins
	if len(plugins) == 0 {
		fmt.Println("There is no plugin")
		return
	}

	mesChan := make(chan gobou.Message)
	// run plugin
	for _, plugInfo := range plugins {
		go func() {
			plug := PluginManager{}
			plug.Load(plugInfo, mesChan)
			plug.Start()
			defer plug.Stop()
			if err := plug.Configure(); err != nil {
				log.Fatalf("Failed to configure Plugin %s", plug.configFile)
				return
			}

			routine := []func() error{
				plug.Main,
				plug.Notify,
				plug.SaveData,
				plug.SaveConfig,
				plug.Wait,
			}
		pluginLoop:
			for {
				for _, f := range routine {
					if err := f(); err != nil {
						log.Fatalln(err)
						break pluginLoop
					}
				}
			}
		}()
	}
	for mes := range mesChan {
		if mes.IsNone() {
			continue
		}
		app.ShowMessage(mes)
		/*
			err := notify.Notify(mes.NotifyMessage.Text)
			if err != nil {
				log.Fatalln(err)
				break
			}
		*/
		time.Sleep(3 * time.Second)
	}
}
