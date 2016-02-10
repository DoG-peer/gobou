package main

import (
	"log"
)

func main() {
	// check arguments
	cli := ParseCliInfo()
	// store pathdata
	app := GetAppPath("gobou")
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
		cli.installInfo.Install(app.PluginDir, app.CacheDir)
		return
	case cli.isGenerate:
		log.Println("generate")
		return
	}

	// search plugins
	pluginPaths, perr := app.GetPlugins()
	if perr != nil {
		log.Fatal(perr)
		return
	}

	// run plugin
	for _, pluginPath := range pluginPaths {
		go func() {
			task := LoadTask(pluginPath, app.PluginConfigDir, app.DataDir)
			task.Start()
			defer task.Stop()
			if err := task.Configure(); err != nil {
				log.Fatalf("Failed to configure Plugin %s", task.configFile)
				return
			}

			for {
				// main task
				if err := task.Main(); err != nil {
					log.Fatalln(err)
					break
				}

				// log
				if err := task.SaveData(); err != nil {
					log.Fatalln(err)
					break
				}

				// change config
				if err := task.SaveConfig(); err != nil {
					log.Fatalln(err)
					break
				}

				// wait
				if err := task.Wait(); err != nil {
					log.Fatalln(err)
					break
				}
			}
		}()
	}
	for {
	}
}
