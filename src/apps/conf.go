package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"runtime"
)

const confFilename string = "apps.yaml"

type Configuration struct {
	QueueUri             string
	User                 string
	Server               string
	ExecutionServers     []string
	SSHPort              string
	RDPPort              string
	Password             string
	WindowsDomain        string
	XMLConfigurationFile string
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
		QueueUri:             "amqp://guest:guest@localhost:5672/",
		SSHPort:              "22",
		RDPPort:              "3389",
		Server:               "62.210.56.76",
		Password:             "ItsPass1942+",
		XMLConfigurationFile: "conf.xml",
		WindowsDomain:        "intra.localdomain.com",
		ExecutionServers:     []string{"62.210.56.76"},
		User:                 "Administrator",
	}
}

func initConf() {
	conf = getDefaultConf()
	usr, err := user.Current()
	if err != nil {
		log.Println(err)
	}
	home := usr.HomeDir
	f := "apps.yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud/apps/"
		err := os.MkdirAll(d, 0755)
		if err == nil {
			f = d + f
		} else {
			log.Println(err)
		}
	}

	if err := ReadMergeConf(&conf, f); err != nil {
		log.Println("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/apps/apps.yaml"
		if err := ReadMergeConf(&conf, alt); err != nil {
			log.Println("No Configuration file found in /etc/nanocloud, using default configuration")
		}
	}
	if err := WriteConf(conf, f); err != nil {
		log.Println(err)
	}
}
