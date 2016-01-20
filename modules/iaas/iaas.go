/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
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
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

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

type VMstatus struct {
	DownloadingVmNames []string
	AvailableVMNames   []string
	BootingVmNames     []string
	RunningVmNames     []string
}

func GetList() (VMstatus, error) {
	var status VMstatus
	var vmIP string

	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/pid/", conf.instDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".pid") {
			continue
		}
		vmIP = strings.Split(file.Name(), "-")[3]
		if checkPort(vmIP, 22) || checkPort(vmIP, 443) || checkPort(vmIP, 3389) {
			status.RunningVmNames = append(status.RunningVmNames, file.Name()[0:len(file.Name())-4])
		} else {
			status.BootingVmNames = append(status.BootingVmNames, file.Name()[0:len(file.Name())-4])
		}
	}

	files, _ = ioutil.ReadDir(fmt.Sprintf("%s/images/", conf.instDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".qcow2") {
			continue
		}
		status.AvailableVMNames = append(status.AvailableVMNames, file.Name()[0:len(file.Name())-6])
	}

	files, _ = ioutil.ReadDir(fmt.Sprintf("%s/downloads/", conf.instDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".qcow2") {
			continue
		}
		status.DownloadingVmNames = append(status.DownloadingVmNames, file.Name()[0:len(file.Name())-6])
	}

	return status, nil
}

func downloadFromUrl(downloadUrl string, dst string) {
	fmt.Println("Downloading", downloadUrl, "to", dst)

	u, err := url.Parse(downloadUrl)
	if err != nil {
		log.Fatal(err)
	}

	splitedPath := strings.Split(u.Path, "/")
	tempDst := filepath.Join(conf.instDir, "downloads", splitedPath[len(splitedPath)-1])
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

func Download(VMName string) error {
	go downloadFromUrl(
		conf.artURL+VMName+".qcow2",
		conf.instDir+"/images/"+VMName+".qcow2")

	return nil
}

func Start(name string) error {
	fmt.Println("Starting : ", name)
	cmd := exec.Command(fmt.Sprintf("%s/scripts/launch-%s.sh", conf.instDir, name))
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start vm, error: %s\n", err)
		return err
	}

	return nil
}

func Stop(name string) error {

	var pidFile string = strings.TrimSpace(fmt.Sprintf("%s/pid/%s.pid\n", conf.instDir, name))
	fmt.Println("stopping : ", name)
	os.Remove(pidFile)
	os.Exit(0)

	return nil
}
