package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
)

const confFilename string = "history.yaml"

type Configuration struct {
	ConnectionString string
	DatabaseName     string
}

var conf Configuration

func ReadMergeConf(out interface{}, filename string) error {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(d, out)
}

func WriteConf(in interface{}, filename string) error {
	d, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, d, 0644)
}

func getDefaultConf() Configuration {
	return Configuration{
		ConnectionString: "history.db",
		DatabaseName:     "bolt",
	}
}

func initConf() {

	conf = getDefaultConf()
	usr, err := user.Current()
	if err != nil {
		log.Error("Failed to get data on current user: ", err)
	}
	home := usr.HomeDir

	f := "history.yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud/history/"
		err := os.MkdirAll(d, 0755)
		if err == nil {
			f = d + f
		} else {
			log.Error("Failed to create necessary directories for config files: ", err)
		}
	}

	if err := ReadMergeConf(&conf, f); err != nil {
		log.Warn("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/history/history.yaml"
		if err := ReadMergeConf(&conf, alt); err != nil {
			log.Warn("No Configuration file found in /etc/nanocloud, using default configuration")
		}

	}
	if err := WriteConf(conf, f); err != nil {
		log.Error("Failed to write configuration file: ", err)
	}
}
