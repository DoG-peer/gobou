package main

import (
	"fmt"
	"log"
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

	// run plugin
	for _, plugInfo := range plugins {
		go func() {
			plug := PluginManager{}
			plug.Load(plugInfo)
			//plug := LoadTask(plugInfo.Path, app.PluginConfigDir, app.DataDir)
			plug.Start()
			defer plug.Stop()
			if err := plug.Configure(); err != nil {
				log.Fatalf("Failed to configure Plugin %s", plug.configFile)
				return
			}

			for {
				// main plug
				if err := plug.Main(); err != nil {
					log.Fatalln(err)
					break
				}

				// log
				if err := plug.SaveData(); err != nil {
					log.Fatalln(err)
					break
				}

				// change config
				if err := plug.SaveConfig(); err != nil {
					log.Fatalln(err)
					break
				}

				// wait
				if err := plug.Wait(); err != nil {
					log.Fatalln(err)
					break
				}
			}
		}()
	}
	for {
	}
}
