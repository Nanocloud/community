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
	"os"
	"strings"

	"github.com/Nanocloud/nano"
)

var module nano.Module

type Configuration struct {
	artURL   string
	instDir  string
	Server   string
	User     string
	SSHPort  string
	Password string
}

type VmInfo struct {
	Ico         string
	Name        string
	DisplayName string
	Status      string
	Locked      bool
}

type hash map[string]interface{}

type VmName struct {
	Name string
}

var conf Configuration

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ListRunningVm(req nano.Request) (*nano.Response, error) {
	var (
		response struct {
			DownloadingVmNames []string
			AvailableVMNames   []string
			BootingVmNames     []string
			RunningVmNames     []string
		}
		vmList      []VmInfo
		icon        string
		locked      bool
		Status      string
		displayName string
	)
	/*	jsonResponse, err := jsonRpcRequest(
			fmt.Sprintf("%s:%s", conf.apiURL, conf.apiPort),
			"Iaas.GetList",
			nil,
		)
		if err != nil {
			module.Log.Error("RPC Call to Iaas API failed: ", err)
			return nil, err
		}*/

	response, err := GetList()
	if err != nil {
		module.Log.Error("Unable to retrieve VM states list: " + err.Error())
		return nano.JSONResponse(500, hash{
			"error": "Unable te retrieve states of VMs",
		}), err
	}

	// TODO: Lots of Data aren't from iaas API
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
		} else if stringInSlice(vmName, response.DownloadingVmNames) {
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
		})
	}

	return nano.JSONResponse(200, vmList), nil
}

func DownloadVm(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"vmname": req.Params["id"],
	}

	err := Download(params["vmname"])
	if err != nil {
		module.Log.Error("Unable to download the specified vm: " + err.Error())
		return nano.JSONResponse(500, hash{
			"error": "Unable to download the specified vm",
		}), err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func StartVm(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"name": req.Params["id"],
	}

	err := Start(params["name"])
	if err != nil {
		module.Log.Error("Error while starting vm: " + err.Error())
		return nano.JSONResponse(500, hash{
			"error": "Unable to start the specified vm",
		}), err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func StopVm(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"Name": req.Params["id"],
	}

	err := Stop(params["Name"])
	if err != nil {
		module.Log.Error("Error while stopping vm: " + err.Error())
		return nano.JSONResponse(500, hash{
			"error": "Unable to stop the specified vm",
		}), err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func env(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
	module = nano.RegisterModule("iaas")

	conf.Server = env("SERVER", "62.210.56.45")
	conf.Password = env("PASSWORD", "ItsPass1942+")
	conf.User = env("USER", "Administrator")
	conf.SSHPort = env("SSH_PORT", "22")

	conf.instDir = os.Getenv("INSTALLATION_DIR")
	if len(conf.instDir) == 0 {
		conf.instDir = "/var/lib/nanocloud"
	}

	conf.artURL = os.Getenv("ARTIFACT_URL")
	if len(conf.artURL) == 0 {
		conf.artURL = "http://releases.nanocloud.org:8080/releases/latest/"
	}

	module.Get("/iaas", ListRunningVm)
	module.Post("/iaas/:id/stop", StopVm)
	module.Post("/iaas/:id/start", StartVm)
	module.Post("/iaas/:id/download", DownloadVm)

	module.Listen()
}
