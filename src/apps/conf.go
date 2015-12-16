package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os/user"
	"path/filepath"
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

func readConfFromPath(path string) error {
	f := filepath.Join(path, confFilename)
	return readMergeConf(&conf, f)
}

func readConfFromHome() error {
	u, err := user.Current()
	if err != nil {
		return err
	}
	path := filepath.Join(u.HomeDir, "/.config/nanocloud")
	return readConfFromPath(path)
}

func initConf() {
	conf = getDefaultConf()
	err := readConfFromHome()
	if err == nil {
		return
	}
	err = readConfFromPath("/etc/nanocloud")
	if err != nil {
		log.Info(confFilename, " is neither found in ~/.config/nanocloud nor in /etc/nanocloud. use default configuration.")
	}
}
