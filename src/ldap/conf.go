package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"

	log "github.com/Sirupsen/logrus"
)

const confFilename string = "ldap.yaml"

type Configuration struct {
	Username  string
	Password  string
	ServerURL string
	QueueURI  string
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
		Username:  "CN=Administrator,CN=Users,DC=intra,DC=localdomain,DC=com",
		Password:  "Nanocloud123+",
		ServerURL: "ldaps://10.20.12.20",
		QueueURI:  "amqp://guest:guest@localhost:5672/",
	}
}

func initConf() {

	conf = getDefaultConf()
	usr, err := user.Current()
	if err != nil {
		log.Error(err)
	}
	home := usr.HomeDir
	f := "ldap.yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud"
		err := os.MkdirAll(d, 0755)
		// creating necessary directories for configuration file if they do not exist

		if err == nil {
			f = d + f
		} else {
			log.Error("Failed to make necessery directories for config files", err)
		}
	}

	// look in ~/.config/nanocloud for config file
	if err := readMergeConf(&conf, f); err != nil {
		log.Warn("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/ldap.yaml"
		// if the config file is not found in ~/.config/nanocloud, look in /etc/nanocloud
		if err := readMergeConf(&conf, alt); err != nil {
			log.Warn("No Configuration file found in /etc/nanocloud, using default configuration")
		}
	}
}
