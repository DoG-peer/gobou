package main

import (
	"log"
	"os/exec"
	"path/filepath"
)

type InstallInfo struct {
	url  string
	name string
}

func (ii *InstallInfo) Install(pluginDir, cacheDir string) {
	gurl, err := parseInstallURL(ii.url)
	src := getSrcPath(ii.name, cacheDir)
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

}

// run git clone command
func clone(gurl, srcPath string) error {
	_, err := exec.Command("git", "clone", gurl, srcPath).Output()
	return err
}

// run go build command
func build(src, dist string) error {
	_, err := exec.Command("go", "build", "-o", dist, src).Output()
	return err
}

func getSrcPath(name, cacheDir string) string {
	return filepath.Join(cacheDir, name)
}

func getDistPath(name, pluginDir string) string {
	return filepath.Join(pluginDir, name)
}

// TODO: if url is illeagal, raise error
func parseInstallURL(url string) (string, error) {
	githubURL := "https://github.com/" + url + ".git"
	return githubURL, nil
}
