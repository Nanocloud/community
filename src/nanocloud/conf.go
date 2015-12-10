package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const confFilename string = "conf.yaml"

type configuration struct {
	RunDir   string
	StagDir  string
	InstDir  string
	Port     string
	FrontDir string
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
		RunDir:   "plugins/running/",
		StagDir:  "plugins/staging/",
		InstDir:  "plugins/installed/",
		FrontDir: "front/",
		Port:     "8080",
	}
}

func initConf() {
	conf = getDefaultConf()

	if err := readMergeConf(&conf, confFilename); err != nil {
		log.Warn("Unable to read/merge the conf file: ", err)
	}
	if err := writeConf(conf, confFilename); err != nil {
		log.Warn("Unable to write the conf file: ", err)
	}
	log.Info("Current conf: ", conf)

	if err := os.MkdirAll(conf.RunDir, 0777); err != nil {
		log.Fatal("Mkdir failed: ", conf.RunDir, err)
	}
	if err := os.MkdirAll(conf.StagDir, 0777); err != nil {
		log.Fatal("Mkdir failed: ", conf.StagDir, err)
	}
	if err := os.MkdirAll(conf.InstDir, 0777); err != nil {
		log.Fatal("Mkdir failed: ", conf.InstDir, err)
	}
}
