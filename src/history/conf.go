package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/user"
	"runtime"
)

const confFilename string = "history.yaml"

type Configuration struct {
	ConnectionString string
	DatabaseName     string
}

var conf Configuration

// Read a configuration file and unmarshal the data in its first parameter
func ReadMergeConf(out interface{}, filename string) error {
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(d, out)
}

// Write the new configuration on the configuration file
func WriteConf(in interface{}, filename string) error {
	d, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, d, 0644)
}

// Default configuration to use if no configuration files are found
func getDefaultConf() Configuration {
	return Configuration{
		ConnectionString: "history.db",
		DatabaseName:     "bolt",
	}
}

func initConf() {

	conf = getDefaultConf()
	usr, err := user.Current()
	if err != nil {
		log.Error("Failed to get data on current user: ", err)
	}
	home := usr.HomeDir

	f := "history.yaml"
	if runtime.GOOS == "linux" {
		d := home + "/.config/nanocloud/history/"
		// Creating necessary directories for configuration file if they do not exist
		err := os.MkdirAll(d, 0755)
		if err == nil {
			f = d + f
		} else {
			log.Error("Failed to create necessary directories for config files: ", err)
		}
	}

	//Look in ~/.config/nanocloud for config file
	if err := ReadMergeConf(&conf, f); err != nil {
		log.Warn("No Configuration file found in ~/.config/nanocloud, now looking in /etc/nanocloud")
		alt := "/etc/nanocloud/history/history.yaml"
		// If the config file is not found in ~/.config/nanocloud, look in /etc/nanocloud
		if err := ReadMergeConf(&conf, alt); err != nil {
			log.Warn("No Configuration file found in /etc/nanocloud, using default configuration")
		}

	}
	// Finally write the fine configuration used in ./config/nanocloud
	if err := WriteConf(conf, f); err != nil {
		log.Error("Failed to write configuration file: ", err)
	}
}
