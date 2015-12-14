package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
)

const confFilename string = "users.yaml"

type Configuration struct {
	DatabaseUri string
	QueueUri    string
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
	f := "users.yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud/"
		err := os.MkdirAll(d, 0755)
		if err == nil {
			f = d + f
		} else {
			log.Println(err)
		}
	}

	if err := ReadMergeConf(&conf, f); err != nil {
		log.Println("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/users.yaml"
		if err := ReadMergeConf(&conf, alt); err != nil {
			log.Println("No Configuration file found in /etc/nanocloud, using default configuration")
		}
	}
	if err := WriteConf(conf, f); err != nil {
		log.Println(err)
	}
}
