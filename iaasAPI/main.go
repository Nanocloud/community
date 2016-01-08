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
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"gopkg.in/yaml.v2"
)

const configFilename string = "conf.yaml"

type Configuration struct {
	InstallationDir string
	ArtifactURL     string
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
	log.Println(in)
	d, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, d, 0644)
}

func getDefaultConf() Configuration {
	return Configuration{
		InstallationDir: "/var/lib/nanocloud",
		ArtifactURL:     "http://releases.nanocloud.org:8080/indiana/",
	}
}

func initConf() {
	conf = getDefaultConf()
	f := "core.yaml"
	if err := ReadMergeConf(&conf, f); err != nil {
		log.Println(err)
	}
	if err := WriteConf(conf, f); err != nil {
		log.Println(err)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func downloadFromUrl(downloadUrl string, dst string) {
	fmt.Println("Downloading", downloadUrl, "to", dst)

	u, err := url.Parse(downloadUrl)
	if err != nil {
		log.Fatal(err)
	}

	tempDst := filepath.Join(conf.InstallationDir, "downloads", u.Path)
	tmpOutput, err := os.Create(tempDst)
	if err != nil {
		fmt.Println("Error while creating", tempDst, "-", err)
		return
	}

	response, err := http.Get(downloadUrl)
	if err != nil {
		fmt.Println("Error while downloading", downloadUrl, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(tmpOutput, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", downloadUrl, "-", err)
		return
	}
	tmpOutput.Close()

	err = os.Rename(tempDst, dst)
	if err != nil {
		fmt.Println("Error while creating", dst, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}

func checkPort(adress string, port int) bool {
	var one []byte
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", adress, port))
	if err != nil {
		return false
	}
	conn.SetReadDeadline(time.Now())
	if _, err := conn.Read(one); err == io.EOF {
		fmt.Printf("Detected closed LAN connection")
		conn.Close()
		conn = nil
		return false
	} else {
		conn.SetReadDeadline(time.Time{})
	}
	defer conn.Close()

	if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
		return false
	}

	return true
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func main() {

	initConf()

	IaasBinDirExists, err := exists(conf.InstallationDir)
	if err != nil || !IaasBinDirExists {
		log.Fatal("You need to run the install binary before running the API")
		os.Exit(1)
	}

	// Setup RPC server
	pRpcServer := rpc.NewServer()
	pRpcServer.RegisterCodec(json.NewCodec(), "application/json")
	pRpcServer.RegisterService(new(Iaas), "")

	http.Handle("/", pRpcServer)

	s := &http.Server{
		Addr:           "0.0.0.0:8082",
		Handler:        pRpcServer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Now listening on http://0.0.0.0:8082")
	log.Fatal(s.ListenAndServe())
}
