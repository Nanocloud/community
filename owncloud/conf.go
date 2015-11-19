package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"runtime"
)

const confFilename string = "conf.yaml"

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
		hostname:      "192.168.1.40/drive",
		protocol:      "https",
	}
}

func initConf() {

	conf = getDefaultConf()
	f := "owncloud.yaml"
	if runtime.GOOS == "linux" {
		d := "/home/antoine/.config/nanocloud/owncloud/"
		err := os.MkdirAll(d, 0644)
		if err == nil {
			f = d + f
		} else {
			log.Println(err)
		}
	}

	if err := ReadMergeConf(&conf, f); err != nil {
		log.Println(err)
	}
	if err := WriteConf(conf, f); err != nil {
		log.Println(err)
	}
}
