package main

import ()

type InstallInfo struct {
	url  string
	name string
}

func (ii *InstallInfo) Install() {

}

func parseInstallUrl(url string) (string, error) {
	githubUrl := "https://github.com/" + url + ".git"
	return githubUrl, nil
}
