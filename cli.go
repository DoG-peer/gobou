package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// CliInfo is
type CliInfo struct {
	isDefault    bool
	isInstall    bool
	isUpdate     bool
	isGenerate   bool
	isConfig     bool
	isHelp       bool
	isTest       bool
	installInfo  InstallInfo
	generateInfo GenerateInfo
	configInfo   Configurator
}

// ShowHelp shows how to use
func (info *CliInfo) ShowHelp() {
	fmt.Println(`how to use:
	gobou
	gobou i user/plugin_name
	gobou i user/plugin_name other_name
	gobou install user/pname
	gobou install user/pname other_name
	gobou u user/plugin_name
	gobou u user/plugin_name other_name
	gobou update user/pname
	gobou update user/pname other_name
	gobou test
	gobou g relative_path
	gobou generate relative_path
	gobou config
	gobou config plugin_name`)
}

// ParseCliInfo transform Args to CliInfo
func ParseCliInfo() CliInfo {
	cinfo := CliInfo{
		isDefault:  false,
		isInstall:  false,
		isUpdate:   false,
		isTest:     false,
		isGenerate: false,
		isHelp:     false,
	}
	if len(os.Args) <= 1 {
		cinfo.isDefault = true
		return cinfo
	}
	switch os.Args[1] {
	case "i":
		fallthrough
	case "install":
		cinfo.isInstall = true
		cinfo.parseInstallInfo(os.Args[2:])
	case "u":
		fallthrough
	case "update":
		cinfo.isUpdate = true
		cinfo.parseInstallInfo(os.Args[2:])
	case "test":
		cinfo.isTest = true
	case "g":
		fallthrough
	case "generate":
		cinfo.isGenerate = true
		cinfo.parseGenerateInfo(os.Args[2:])
	case "config":
		cinfo.isConfig = true
		cinfo.parseConfigInfo(os.Args[2:])
	default:
		cinfo.isHelp = true
	}
	return cinfo
}

/*
	gobou i user/pname
	gobou i user/pname other_name
	gobou install user/pname
	gobou install user/pname other_name
	gobou u user/pname
	gobou u user/pname other_name
	gobou update user/pname
	gobou update user/pname other_name
*/
func (info *CliInfo) parseInstallInfo(args []string) {
	switch len(args) {
	case 1:
		info.installInfo = InstallInfo{
			url:  args[0],
			name: filepath.Base(args[0]),
		}
	case 2:
		info.installInfo = InstallInfo{
			url:  args[0],
			name: args[1],
		}
	default:
		info.isInstall = false
		info.isUpdate = false
		info.isHelp = true
	}

}

func (info *CliInfo) parseGenerateInfo(args []string) {
	if len(args) == 1 {
		info.generateInfo = GenerateInfo{
			path: args[0],
			name: filepath.Base(args[0]),
		}
	} else {
		info.isGenerate = false
		info.isHelp = true
	}

}

func (info *CliInfo) parseConfigInfo(args []string) {
	if len(args) == 1 {
		info.configInfo = Configurator{
			isMain: false,
			plugin: args[0],
		}
	} else {
		info.configInfo = Configurator{
			isMain: true,
			plugin: "",
		}
	}
}
