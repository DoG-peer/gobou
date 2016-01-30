package main

import (
	"log"
	"time"
)

func main() {
	app := GetAppPath("gobou")
	if err := app.PrepareDirs(); err != nil {
		log.Fatal(err)
		return
	}
	pluginPaths, perr := app.GetPlugins()
	if perr != nil {
		log.Fatal(perr)
		return
	}

	for _, pluginPath := range pluginPaths {
		go func() {
			task := LoadTask(pluginPath, app.PluginConfigDir, app.DataDir)
			task.Start()
			defer task.Stop()
			if err := task.Configure(); err != nil {
				log.Fatalf("Failed to configurePlugin %s", task.configFile)
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
	time.Sleep(10 * time.Second)

}
