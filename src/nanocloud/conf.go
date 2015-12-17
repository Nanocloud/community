package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const confFilename string = "nanocloud.yaml"

type configuration struct {
	RunDir      string
	StagDir     string
	InstDir     string
	Port        string
	FrontDir    string
	UploadDir   string
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
		UploadDir:   "uploads/",
		Port:        "8080",
		DatabaseUri: "postgres://nanocloud@localhost/nanocloud?sslmode=disable",
		QueueUri:    "amqp://guest:guest@localhost:5672/",
	}
}

func readConfFromPath(path string) error {
	f := filepath.Join(path, confFilename)
	err := readMergeConf(&conf, f)
	if !os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"module": moduleName,
			"error":  err,
		}).Error("Unable to read from the configuration file.")
	}
	return err
}

func readConfFromHome() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	path := filepath.Join(u.HomeDir, ".config/nanocloud")
	return readConfFromPath(path)
}

// Trying at first to read and merge from home, then from etc.
// If both failed then keep the default.
func initConf() {
	conf = getDefaultConf()
	readConfFromHome()
	readConfFromPath("/etc/nanocloud")
	log.WithFields(log.Fields{
		"module": moduleName,
		"conf":   conf,
	}).Info("Configuration used")
}
