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
	RunDir  string
	StagDir string
	InstDir string
	Port    string
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
	log.Println(in)
	d, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, d, 0644)
}

func getDefaultConf() Configuration {
	return Configuration{
		RunDir:  "plugins/running/",
		StagDir: "plugins/staging/",
		InstDir: "plugins/installed/",
		Port:    "8080",
	}
}

func initConf() {

	conf = getDefaultConf()
	f := "core.yaml"
	if runtime.GOOS == "linux" {
		d := "/home/antoine/nanocloud/core/"
		err := os.MkdirAll(d, 0777)
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
