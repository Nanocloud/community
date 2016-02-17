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

package iaas

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	VMShutdownFailed = errors.New("VM Shutdown Failed")
	VMStartupFailed  = errors.New("VM Startup Failed")
	VMDownloadFailed = errors.New("VM Download Failed")
)

type Iaas struct {
	Server   string
	Password string
	User     string
	SSHPort  string
	InstDir  string
	ArtURL   string
}

type DownloadingVm struct {
	Name        string `json:"name"`
	CurrentSize string `json:"current_size"`
	TotalSize   string `json:"total_size"`
}

type VMstatus struct {
	DownloadingVmNames []DownloadingVm
	AvailableVMNames   []string
	BootingVmNames     []string
	RunningVmNames     []string
}

type VmInfo struct {
	Ico         string
	Name        string
	DisplayName string
	Status      string
	Locked      bool
	CurrentSize string
	TotalSize   string
}

func New(Server, Password, User, SSHPort, InstDir, ArtURL string) *Iaas {
	return &Iaas{
		Server:   Server,
		Password: Password,
		User:     User,
		SSHPort:  SSHPort,
		InstDir:  InstDir,
		ArtURL:   ArtURL,
	}
}

func (i *Iaas) stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func vmIsDownloading(vmName string, tab []DownloadingVm) bool {
	for _, val := range tab {
		if vmName == val.Name {
			return true
		}
	}
	return false
}

func getCurrentSize(vmName string, tab []DownloadingVm) string {
	for _, val := range tab {
		if vmName == val.Name {
			return val.CurrentSize
		}
	}
	return ""
}

func getTotalSize(vmName string, tab []DownloadingVm) string {
	for _, val := range tab {
		if vmName == val.Name {
			return val.TotalSize
		}
	}
	return ""
}

func (i *Iaas) CheckVMStates(response VMstatus) []VmInfo {
	var (
		locked      bool
		icon        string
		Status      string
		displayName string
		vmList      []VmInfo
	)

	for _, vmName := range response.AvailableVMNames {

		locked = false
		if strings.Contains(vmName, "windows") {
			if strings.Contains(vmName, "winapps") {
				icon = "settings_applications"
				displayName = "Execution environment"
			} else {
				icon = "windows"
				displayName = "Windows Active Directory"
			}
		} else {
			if strings.Contains(vmName, "drive") {
				icon = "storage"
				displayName = "Drive"
			} else if strings.Contains(vmName, "licence") {
				icon = "vpn_lock"
				displayName = "Windows Licence service"
			} else {
				icon = "apps"
				locked = true
				displayName = "Haptic"
			}
		}

		if i.stringInSlice(vmName, response.RunningVmNames) {
			Status = "running"
		} else if i.stringInSlice(vmName, response.BootingVmNames) {
			Status = "booting"
		} else if vmIsDownloading(vmName, response.DownloadingVmNames) {
			Status = "download"
		} else if i.stringInSlice(vmName, response.AvailableVMNames) {
			Status = "available"
		}
		vmList = append(vmList, VmInfo{
			Ico:         icon,
			Name:        vmName,
			DisplayName: displayName,
			Status:      Status,
			Locked:      locked,
			CurrentSize: getCurrentSize(vmName, response.DownloadingVmNames),
			TotalSize:   getTotalSize(vmName, response.DownloadingVmNames),
		})
	}
	return vmList
}

func (i *Iaas) CheckRDS() bool {

	cmd := exec.Command(
		"sshpass", "-p", i.Password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "ConnectTimeout=1",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", i.SSHPort,
		fmt.Sprintf(
			"%s@%s",
			i.User,
			i.Server,
		),
		"powershell.exe \"Write-Host (Get-Service -Name RDMS).status\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to check windows' state", err, string(response))
		return false
	}

	if strings.Contains(string(response), "Running") {
		return true
	}
	return false
}

func generateDownloadURL(url string, vm string) string {
	return url + vm
}

func (i *Iaas) GetList() (VMstatus, error) {
	var status VMstatus
	running := i.CheckRDS()
	if running {
		status.AvailableVMNames = append(status.AvailableVMNames, "windows")
		status.RunningVmNames = append(status.RunningVmNames, "windows")
		return status, nil
	}
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/pid/", i.InstDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".pid") {
			continue
		}
		running := i.CheckRDS()
		if running {
			status.RunningVmNames = append(status.RunningVmNames, file.Name()[0:len(file.Name())-4])
		} else {
			status.BootingVmNames = append(status.BootingVmNames, file.Name()[0:len(file.Name())-4])
		}
	}

	files, _ = ioutil.ReadDir(fmt.Sprintf("%s/images/", i.InstDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".qcow2") {
			continue
		}
		status.AvailableVMNames = append(status.AvailableVMNames, file.Name()[0:len(file.Name())-6])
	}

	files, _ = ioutil.ReadDir(fmt.Sprintf("%s/downloads/", i.InstDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".qcow2") {
			continue
		}
		fi, err := os.Open(filepath.Join(i.InstDir, "downloads", fileName))
		if err != nil {
			log.Error("Couldn't open downloading file: ", err)
			continue
		}
		fiStat, err := fi.Stat()
		if err != nil {
			log.Error("Couldn't stat downloading file: ", err)
			continue
		}
		response, err := http.Head(generateDownloadURL(i.ArtURL, fileName))
		if err != nil {
			log.Error("Error while checking file size ", err)
			continue
		}
		status.AvailableVMNames = append(status.AvailableVMNames, file.Name()[0:len(file.Name())-6])
		status.DownloadingVmNames = append(status.DownloadingVmNames, DownloadingVm{
			Name:        file.Name()[0 : len(file.Name())-6],
			CurrentSize: strconv.FormatInt(fiStat.Size(), 10),
			TotalSize:   strconv.FormatInt(response.ContentLength, 10),
		})
	}

	return status, nil
}

func (i *Iaas) Stop(name string) error {
	log.Info("stopping : ", name)

	cmd := exec.Command(
		"sshpass", "-p", i.Password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", i.SSHPort,
		fmt.Sprintf(
			"%s@%s",
			i.User,
			i.Server,
		),
		"powershell.exe \"Stop-Computer -Force\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to execute sshpass command to shutdown windows", err, string(response))
		return VMShutdownFailed
	}

	return nil
}

func (i *Iaas) Start(name string) error {
	log.Info("Starting : ", name)
	cmd := exec.Command(fmt.Sprintf("%s/scripts/launch-%s.sh", i.InstDir, name))
	err := cmd.Start()
	if err != nil {
		log.Error("Failed to start vm: ", err)
		return VMStartupFailed
	}
	return nil
}

func (i *Iaas) downloadFromUrl(downloadUrl string, dst string) error {
	log.Info("Downloading ", downloadUrl, "to ", dst)

	u, err := url.Parse(downloadUrl)
	if err != nil {
		log.Error("Couldn't parse the VM's URL: ", err)
		return VMDownloadFailed
	}

	splitedPath := strings.Split(u.Path, "/")
	tempDst := filepath.Join(i.InstDir, "downloads", splitedPath[len(splitedPath)-1])
	tmpOutput, err := os.Create(tempDst)
	if err != nil {
		log.Error("Error while creating", tempDst, "-", err)
		return VMDownloadFailed
	}

	response, err := http.Get(downloadUrl)
	if err != nil {
		log.Error("Error while downloading", downloadUrl, "-", err)
		return VMDownloadFailed
	}
	defer response.Body.Close()

	n, err := io.Copy(tmpOutput, response.Body)
	if err != nil {
		log.Error("Error while downloading", downloadUrl, "-", err)
		return VMDownloadFailed
	}
	tmpOutput.Close()

	err = os.Rename(tempDst, dst)
	if err != nil {
		log.Error("Error while creating", dst, "-", err)
		return VMDownloadFailed
	}

	log.Info(n, "bytes downloaded.")
	return nil
}

func (i *Iaas) Download(VMName string) {
	i.downloadFromUrl(
		generateDownloadURL(i.ArtURL, VMName)+".qcow2",
		i.InstDir+"/images/"+VMName+".qcow2")
}
