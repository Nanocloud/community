package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

// Read a configuration file and unmarshal the data in its first parameter
func readMergeConf(out interface{}, filename string) error {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(d, out)
}

// Write the new configuration on the configuration file
func writeConf(in interface{}, filename string) error {
	d, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, d, 0644)
}

// Default configuration to use if no configuration files are found
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
		// creating necessary directories for configuration file if they do not exist
		if err == nil {
			f = d + f
		} else {
			log.Println(err)
		}
	}

	// look in ~/.config/nanocloud for config file
	if err := readMergeConf(&conf, f); err != nil {
		log.Warn("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/apps/apps.yaml"
		// if the config file is not found in ~/.config/nanocloud, look in /etc/nanocloud
		if err := readMergeConf(&conf, alt); err != nil {
			log.Warn("No Configuration file found in /etc/nanocloud, using default configuration")
		}
	}
	// finally write the final configuration used in ./config/nanocloud
	if err := writeConf(conf, f); err != nil {
		log.Error("Failed to write configuration file for plugin apps: ", err)
	}
}
