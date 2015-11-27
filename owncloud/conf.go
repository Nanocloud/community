package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
)

const confFilename string = "owncloud.yaml"

type Configuration struct {
	adminLogin    string
	adminPassword string
	protocol      string
	hostname      string
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
		adminLogin:    "drive_admin",
		adminPassword: "BJboVHDiawECoDt",
		hostname:      "192.168.1.39/drive",
		protocol:      "https",
	}
}

func initConf() {

	conf = getDefaultConf()
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	home := usr.HomeDir
	f := "owncloud.yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud/owncloud/"
		err := os.MkdirAll(d, 0755)
		if err == nil {
			f = d + f
		} else {
			log.Println(err)
		}
	}

	if err := ReadMergeConf(&conf, f); err != nil {
		log.Println("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/owncloud/owncloud.yaml"
		if err := ReadMergeConf(&conf, alt); err != nil {
			log.Println("No Configuration file found in /etc/nanocloud, using default configuration")
		}

	}
	if err := WriteConf(conf, f); err != nil {
		log.Println(err)
	}
}
