package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// InstallInfo is
type InstallInfo struct {
	url  string
	name string
}

// Install runs by gobou install hoge/foo
func (ii *InstallInfo) Install(pluginDir, cacheDir, pluginConfigDir string) {
	gurl, err := parseInstallURL(ii.url)
	srcDir := filepath.Join(cacheDir, "src")
	src := getSrcPath(ii.name, srcDir)
	dist := getDistPath(ii.name, pluginDir)
	if err != nil {
		return
	}
	err = clone(gurl, src)
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = build(src, dist)
	if err != nil {
		log.Fatalln(err)
		return
	}

	err = initConfigFile(filepath.Join(pluginConfigDir, ii.name))
	if err != nil {
		log.Fatalln(err)
		return
	}

}

// run git clone command
func clone(gurl, srcPath string) error {
	out, err := exec.Command("git", "clone", gurl, srcPath).CombinedOutput()
	log.Println(string(out))
	return err
}

// run go build command
func build(src, dist string) error {
	pwd, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	relsrc, err := filepath.Rel(pwd, src)
	if err != nil {
		return err
	}
	out, err := exec.Command("go", "build", "-o", dist, relsrc).CombinedOutput()
	log.Println(string(out))
	return err
}

func getSrcPath(name, srcDir string) string {
	return filepath.Join(srcDir, name)
}

func getDistPath(name, pluginDir string) string {
	return filepath.Join(pluginDir, name)
}

// TODO: if url is illeagal, raise error
func parseInstallURL(url string) (string, error) {
	githubURL := "https://github.com/" + url + ".git"
	return githubURL, nil
}

func initConfigFile(file string) error {
	_, e := os.Stat(file)
	if os.IsNotExist(e) {
		return ioutil.WriteFile(file, []byte("{}"), os.ModePerm)
	}
	return nil
}
