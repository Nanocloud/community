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
	ScriptsDir string
	Username   string
	Password   string
	ServerURL  string
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
		ScriptsDir: "/home/antoine/protov2/ldap/",
		Username:   "CN=Administrator,CN=Users,DC=intra,DC=localdomain,DC=com",
		Password:   "Nanocloud123+",
		ServerURL:  "ldaps://10.20.12.20",
	}
}

func initConf() {

	conf = getDefaultConf()
	f := "ldap.yaml"
	if runtime.GOOS == "linux" {
		d := "/home/antoine/.config/nanocloud/ldap/"
		err := os.MkdirAll(d, 0644)
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
