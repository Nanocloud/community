/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	log "github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
		DatabaseUri: "postgres://nanocloud@localhost/nanocloud?sslmode=disable",
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
	if runtime.GOOS != "windows" {
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
}
