package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
)

const confFilename string = "iaas.yaml"

type configuration struct {
	Url  string
	Port string
}

var conf configuration

func readMergeConf(out interface{}, filename string) error {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(d, out)
}

func writeConf(in interface{}, filename string) error {
	d, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, d, 0644)
}

func getDefaultConf() configuration {
	return configuration{
		Url:  "http://192.168.1.40",
		Port: "8082",
	}
}

func readConfFromPath(path string) error {
	f := filepath.Join(path, confFilename)
	return readMergeConf(&conf, f)
}

func readConfFromHome() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	path := filepath.Join(u.HomeDir, "/.config/nanocloud")
	return readConfFromPath(path)
}

func initConf() {
	conf = getDefaultConf()
	err := readConfFromHome()
	if err == nil {
		return
	}
	err = readConfFromPath(filepath.Join("/etc/nanocloud", confFilename))
	if err != nil {
		log.Info(confFilename, " is neither found in ~/.config/nanocloud nor in /etc/nanocloud. using default configuration.")
	}
}
