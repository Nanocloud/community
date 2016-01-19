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
	"fmt"
	"github.com/Nanocloud/nano"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var module nano.Module

var apiURL string
var apiPort string

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
			Result struct {
				DownloadingVmNames []string
				AvailableVMNames   []string
				BootingVmNames     []string
				RunningVmNames     []string
			}
			Error string
			Id    int
		}
		vmList      []VmInfo
		icon        string
		locked      bool
		Status      string
		displayName string
	)
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", apiURL, apiPort),
		"Iaas.GetList",
		nil,
	)
	if err != nil {
		module.Log.Error("RPC Call to Iaas API failed: ", err)
		return nil, err
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		module.Log.Error("Failed to Unmarshal response from Iaas API: ", err)
		return nil, err
	}

	// TODO: Lots of Data aren't from iaas API
	for _, vmName := range response.Result.AvailableVMNames {

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

		if stringInSlice(vmName, response.Result.RunningVmNames) {
			Status = "running"
		} else if stringInSlice(vmName, response.Result.BootingVmNames) {
			Status = "booting"
		} else if stringInSlice(vmName, response.Result.DownloadingVmNames) {
			Status = "download"
		} else if stringInSlice(vmName, response.Result.AvailableVMNames) {
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

	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", apiURL, apiPort),
		"Iaas.Download",
		params,
	)
	if err != nil {
		module.Log.Error("RPC Call to Iaas API failed: ", err)
		return nil, err
	}

	return &nano.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Body:        jsonResponse,
	}, nil
}

func StartVm(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"name": req.Params["id"],
	}
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", apiURL, apiPort),
		"Iaas.Start",
		params,
	)
	if err != nil {
		module.Log.Error("RPC Call to Iaas API failed: ", err)
		return nil, err
	}

	return &nano.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Body:        jsonResponse,
	}, nil
}

func StopVm(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"Name": req.Params["id"],
	}
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", apiURL, apiPort),
		"Iaas.Stop",
		params,
	)
	if err != nil {
		module.Log.Error("RPC Call to Iaas API failed: ", err)
		return nil, err
	}
	return &nano.Response{
		StatusCode:  200,
		ContentType: "application/json",
		Body:        jsonResponse,
	}, nil
}

func jsonRpcRequest(url string, method string, param map[string]string) ([]byte, error) {
	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      1,
		"params":  []map[string]string{0: param},
	})
	if err != nil {
		module.Log.Errorf("Marshal: %v", err)
		return nil, err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		module.Log.Errorf("Post: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		module.Log.Errorf("ReadAll: %v", err)
		return nil, err
	}

	return body, nil
}

func main() {
	module = nano.RegisterModule("iaas")

	apiURL = os.Getenv("API_URL")
	if len(apiURL) == 0 {
		apiURL = "http://192.168.1.40"
	}

	apiPort = os.Getenv("API_PORT")
	if len(apiPort) == 0 {
		apiPort = "8082"
	}

	module.Get("/iaas", ListRunningVm)
	module.Post("/iaas/:id", StopVm)
	module.Post("/iaas/:id/start", StartVm)
	module.Post("/iaas/:id/download", DownloadVm)

	module.Listen()
}
