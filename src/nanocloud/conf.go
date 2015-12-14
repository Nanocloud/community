package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
)

const confFilename string = "nanocloud.yaml"

type configuration struct {
	RunDir      string
	StagDir     string
	InstDir     string
	Port        string
	FrontDir    string
	DatabaseUri string
	QueueUri    string
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
		RunDir:      "plugins/running/",
		StagDir:     "plugins/staging/",
		InstDir:     "plugins/installed/",
		FrontDir:    "front/",
		Port:        "8080",
		DatabaseUri: "postgres://localhost/nanocloud?sslmode=disable",
		QueueUri:    "amqp://guest:guest@localhost:5672/",
	}
}

func initConf() {
	conf = getDefaultConf()
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	home := usr.HomeDir
	f := "nanocloud.yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud/"
		err := os.MkdirAll(d, 0755)
		if err == nil {
			f = d + f
		} else {
			log.Println(err)
		}
	}

	if err := readMergeConf(&conf, f); err != nil {
		log.Println("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/nanocloud.yaml"
		if err := readMergeConf(&conf, alt); err != nil {
			log.Println("No Configuration file found in /etc/nanocloud, using default configuration")
		}
	}
}
