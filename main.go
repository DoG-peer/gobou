package main

import (
	"github.com/dullgiulio/pingo"
	"log"
	"path/filepath"
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
			task.configFile = filepath.Join(app.PluginConfigDir, task.name)
			task.dataFile = filepath.Join(app.DataDir, task.name)

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
