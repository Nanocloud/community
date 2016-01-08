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
	"io/ioutil"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

var ()

type Iaas struct{}

type NoArgs struct{}
type VMName struct {
	Name string
}

type GetIaasListReply struct {
	DownloadingVmNames []string
	AvailableVMNames   []string
	BootingVmNames     []string
	RunningVmNames     []string
}

func (p *Iaas) GetList(r *http.Request, args *NoArgs, reply *GetIaasListReply) error {
	var vmIP string

	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/pid/", conf.InstallationDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".pid") {
			continue
		}
		vmIP = strings.Split(file.Name(), "-")[3]
		if checkPort(vmIP, 22) || checkPort(vmIP, 443) || checkPort(vmIP, 3389) {
			reply.RunningVmNames = append(reply.RunningVmNames, file.Name()[0:len(file.Name())-4])
		} else {
			reply.BootingVmNames = append(reply.BootingVmNames, file.Name()[0:len(file.Name())-4])
		}
	}

	files, _ = ioutil.ReadDir(fmt.Sprintf("%s/images/", conf.InstallationDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".qcow2") {
			continue
		}
		reply.AvailableVMNames = append(reply.AvailableVMNames, file.Name()[0:len(file.Name())-6])
	}

	files, _ = ioutil.ReadDir(fmt.Sprintf("%s/downloads/", conf.InstallationDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".qcow2") {
			continue
		}
		reply.DownloadingVmNames = append(reply.DownloadingVmNames, file.Name()[0:len(file.Name())-6])
	}

	return nil
}

type StatusReply struct {
	status string
}

func (p *Iaas) GetStatus(r *http.Request, args *VMName, reply *StatusReply) error {
	var vmIP string = strings.Split(args.Name, "-")[3]

	if checkPort(vmIP, 443) {
		reply.status = "running"
	} else if Exists(fmt.Sprintf("%s/pid/%s.pid", conf.InstallationDir, args.Name)) {
		reply.status = "booting"
	} else if Exists(fmt.Sprintf("%s/images/%s.qcow2", conf.InstallationDir, args.Name)) {
		reply.status = "available"
	} else if Exists(fmt.Sprintf("%s/downloads/%s.qcow2", conf.InstallationDir, args.Name)) {
		reply.status = "downloading"
	} else {
		reply.status = "unknown"
	}

	return nil
}

type BoolReply struct {
	Success bool
}

func (p *Iaas) Stop(r *http.Request, args *VMName, reply *BoolReply) error {
	var (
		socketFile string = strings.TrimSpace(fmt.Sprintf("%s/sockets/%s.socket\n", conf.InstallationDir, args.Name))
	)
	fmt.Printf("Shut down VM: «%s»\n", socketFile)

	buf := make([]byte, 1024)
	connection, err := net.Dial("unix", socketFile)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer connection.Close()

	_, err = connection.Read(buf[:])
	if err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond)
	_, err = connection.Write([]byte("system_powerdown\n"))
	if err != nil {
		fmt.Println("Write error:", err)
	}

	time.Sleep(1000 * time.Millisecond)
	_, err = connection.Read(buf[:])
	if err != nil {
		return err
	}

	reply.Success = true
	return nil
}

func (p *Iaas) Start(r *http.Request, args *VMName, reply *BoolReply) error {

	fmt.Println("Starting : ", args)
	cmd := exec.Command("nohup", fmt.Sprintf("%s/scripts/launch-%s.sh", conf.InstallationDir, args.Name), "&")
	err := cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start vm, error: %s\n", err)
		reply.Success = false
		return err
	}

	reply.Success = true
	return nil
}

type ImageArgs struct {
	VMName string
}

func (p *Iaas) Download(r *http.Request, args *ImageArgs, reply *BoolReply) error {
	go downloadFromUrl(
		conf.ArtifactURL+args.VMName+".qcow2",
		conf.InstallationDir+"/images/"+args.VMName+".qcow2")

	reply.Success = true

	return nil
}
