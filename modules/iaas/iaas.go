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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	VMShutdownFailed = errors.New("VM Shutdown Failed")
	VMStartupFailed  = errors.New("VM Startup Failed")
	VMDownloadFailed = errors.New("VM Download Failed")
	VMNotFound       = errors.New("Specified VM does not exists")
)

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
	Ico         string `json:"ico"`
	Name        string `json:"-"`
	DisplayName string `json:"display_name"`
	Status      string `json:"status"`
	Locked      bool   `json:"locked"`
	CurrentSize string `json:"current_size"`
	TotalSize   string `json:"total_size"`
}

func stringInSlice(a string, list []string) bool {
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

func CheckVMStates(response VMstatus) []VmInfo {
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

		if stringInSlice(vmName, response.RunningVmNames) {
			Status = "running"
		} else if stringInSlice(vmName, response.BootingVmNames) {
			Status = "booting"
		} else if vmIsDownloading(vmName, response.DownloadingVmNames) {
			Status = "download"
		} else if stringInSlice(vmName, response.AvailableVMNames) {
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

func CheckRDS() bool {
	resp, err := http.Get("http://" + conf.Server + ":9090/checkrds")
	if err != nil {
		log.Error(err)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		return false
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return false
	}

	if strings.Contains(string(b), "Running") {
		return true
	}
	return false
}

func generateDownloadURL(url string, vm string) string {
	return url + vm
}

func GetList() (VMstatus, error) {
	var status VMstatus
	running := CheckRDS()
	if running {
		status.AvailableVMNames = append(status.AvailableVMNames, "windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64")
		status.RunningVmNames = append(status.RunningVmNames, "windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64")
		return status, nil
	}
	files, _ := ioutil.ReadDir(fmt.Sprintf("%s/pid/", conf.instDir))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, ".pid") {
			continue
		}
		running := CheckRDS()
		if running {
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
		fi, err := os.Open(filepath.Join(conf.instDir, "downloads", fileName))
		if err != nil {
			log.Error("Couldn't open downloading file: ", err)
			continue
		}
		fiStat, err := fi.Stat()
		if err != nil {
			log.Error("Couldn't stat downloading file: ", err)
			continue
		}
		response, err := http.Head(generateDownloadURL(conf.artURL, fileName))
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

func Stop(name string) error {
	log.Info("stopping : ", name)
	resp, err := http.Get("http://" + conf.Server + ":9090/shutdown")
	if err != nil {
		log.Error(err)
		return VMShutdownFailed
	}

	if resp.StatusCode != http.StatusOK {
		return VMShutdownFailed
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return VMShutdownFailed
	}

	var ret map[string]map[string]bool
	err = json.Unmarshal(b, &ret)
	if ret["data"]["success"] == false {
		return VMShutdownFailed
	}
	return nil
}

func createQcow() error {
	cmd := exec.Command(
		"qemu-img",
		"create",
		"-f",
		"qcow2",
		"-o",
		"preallocation=metadata",
		conf.instDir+"/downloads/win.qcow2",
		"30G",
	)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(resp))
		log.Error(err)
		return err
	}
	return nil
}

func createIso() error {
	plazaLocation := os.Getenv("PLAZA_LOCATION")

	if len(plazaLocation) == 0 {
		return errors.New("plaza cannot be found")
	}

	disk := path.Join(conf.root, "disk")

	err := copyFile(plazaLocation, path.Join(disk, "plaza.exe"))
	if err != nil {
		return err
	}

	cmd := exec.Command(
		"genisoimage",
		"-o",
		conf.instDir+"/downloads/autoplaza.iso",
		"-J",
		"-r",
		disk,
	)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(string(resp))
		log.Error(err)
		return err
	}
	return nil
}

func downloadIso() error {
	tab := strings.Split(conf.windowsURL, "/")
	if _, err := os.Stat(conf.instDir + "/downloads/" + tab[len(tab)-1]); err == nil {
		return nil
	}
	out, err := os.Create(conf.instDir + "/downloads/" + tab[len(tab)-1])
	if err != nil {
		log.Error(err)
		return err
	}
	defer out.Close()
	resp, err := http.Get(conf.windowsURL)
	if err != nil {
		log.Error(err)
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func bootWindows() error {
	tab := strings.Split(conf.windowsURL, "/")
	cmd := exec.Command(
		"qemu-system-x86_64",
		"-m", "4096",
		"-cpu", "host",
		"-machine", "accel=kvm",
		"-smp", "4",
		"-vnc", ":2",
		"-device", "virtio-net,netdev=user.0",
		"-boot", "once=d",
		"-machine", "type=pc,accel=kvm",
		"-drive", "file="+conf.instDir+"/downloads/autoplaza.iso"+",index=0,media=cdrom",
		"-drive", "file="+conf.instDir+"/downloads/"+tab[len(tab)-1]+",index=1,media=cdrom",
		"-drive", "file="+conf.instDir+"/downloads/win.qcow2"+",index=2,if=virtio,cache=writeback,discard=ignore",
		"-netdev", "user,id=user.0",
		"-vga", "qxl",
		"-global", "qxl-vga.vram_size=33554432",
	)
	resp, err := cmd.CombinedOutput()
	if err != nil {
		log.Error(err)
		log.Error(string(resp))
		return err
	}
	err = os.Rename(conf.instDir+"/downloads/win.qcow2", conf.instDir+"/images/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.qcow2")
	if err != nil {
		log.Error(err)
	}
	return nil
}

func Create() error {
	err := downloadIso()
	if err != nil {
		return err
	}
	log.Error("CREATING QCOW2")
	err = createQcow()
	if err != nil {
		return err
	}
	log.Error("CREATING ISO")
	err = createIso()
	if err != nil {
		return err
	}
	log.Error("BOOTING WINDOWS")
	err = bootWindows()
	if err != nil {
		return err
	}
	return nil
}

func Start(name string) error {
	log.Info("Starting : ", name)
	_, err := os.Stat(fmt.Sprintf("%s/images/%s.qcow2", conf.instDir, name))
	if os.IsNotExist(err) {
		log.Error("Can't find ", name)
		return VMNotFound
	}
	cmd := exec.Command(fmt.Sprintf("%s/scripts/launch-%s.sh", conf.root, name))
	err = cmd.Start()
	if err != nil {
		log.Error("Failed to start vm: ", err)
		return VMStartupFailed
	}
	return nil
}

func downloadFromUrl(downloadUrl string, dst string) error {
	log.Info("Downloading ", downloadUrl, "to ", dst)

	u, err := url.Parse(downloadUrl)
	if err != nil {
		log.Error("Couldn't parse the VM's URL: ", err)
		return VMDownloadFailed
	}

	splitedPath := strings.Split(u.Path, "/")
	tempDst := filepath.Join(conf.instDir, "downloads", splitedPath[len(splitedPath)-1])
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
		// If copy fails that's probably because downloads and image are not on the same partition
		// In this case we should make an hard copy
		err = copyFile(tempDst, dst)
		if err != nil {
			log.Error("Error while creating ", dst, "-", err)
			return VMDownloadFailed
		}
		err = os.Remove(tempDst)
		if err != nil {
			log.Error("Error while removing ", tempDst, "-", err)
			return VMDownloadFailed
		}

	}

	log.Info(n, "bytes downloaded.")
	return nil
}

func Download(VMName string) {
	downloadFromUrl(
		generateDownloadURL(conf.artURL, VMName)+".qcow2",
		conf.instDir+"/images/"+VMName+".qcow2")
}
