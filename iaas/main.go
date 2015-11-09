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
	"io/ioutil"
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"

	"os"
	"strings"

	"github.com/natefinch/pie"
)

var (
	name = "iaas"
	srv  pie.Server
)

type VmInfo struct {
	Ico         string
	Name        string
	DisplayName string
	Status      string
	Locked      bool
}

type api struct{}

type VmName struct {
	Name string
}

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ListRunningVm(args PlugRequest, reply *PlugRequest) error {

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
		status      string
		displayName string
	)
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.GetList",
		nil,
	)
	if err != nil {
		log.Println(err) // for on-screen debug output
		return nil
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		log.Println(err) // for on-screen debug output
		return nil
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
			status = "running"
		} else if stringInSlice(vmName, response.Result.BootingVmNames) {
			status = "booting"
		} else if stringInSlice(vmName, response.Result.DownloadingVmNames) {
			status = "download"
		} else if stringInSlice(vmName, response.Result.AvailableVMNames) {
			status = "available"
		}
		vmList = append(vmList, VmInfo{
			Ico:         icon,
			Name:        vmName,
			DisplayName: displayName,
			Status:      status,
			Locked:      locked,
		})
	}

	jsonOutput, _ := json.Marshal(vmList)
	reply.Body = string(jsonOutput)
	return err
}

func DownloadVm(args PlugRequest, reply *PlugRequest) {
	var t VmName
	err := json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
	}

	var (
		params = map[string]string{
			"vmname": t.Name,
		}
		response struct {
			Result struct {
				Success bool
			}
		}
	)

	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.Download",
		params,
	)
	if err != nil {
		log.Println(err) // for on-screen debug output
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		log.Println(err) // for on-screen debug output
	}

}

func DownloadStatus(args PlugRequest, reply *PlugRequest) {
	var (
		response struct {
			Result struct {
				AvailableVMNames   []string
				RunningVmNames     []string
				DownloadInProgress bool
			}
			Error string
			Id    int
		}
	)
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.GetList",
		nil,
	)
	if err != nil {
		log.Println(err) // for on-screen debug output
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		log.Println(err) // for on-screen debug output
	}

}

func StartVm(args PlugRequest, reply *PlugRequest) {
	var t VmName
	err := json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
	}

	var (
		params = map[string]string{
			"name": t.Name,
		}
		response struct {
			Result struct {
				Success bool
			}
		}
	)

	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.Start",
		params,
	)
	if err != nil {
		log.Println(err)
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		log.Println(err)
	}

}

func StopVm(args PlugRequest, reply *PlugRequest) {
	var t VmName
	err := json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
	}
	var params = map[string]string{
		"Name": t.Name,
	}
	_, err = jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.Stop",
		params,
	)
	if err != nil {
		log.Println(err)
	}
}
func jsonRpcRequest(url string, method string, param map[string]string) (string, error) {

	data, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      1,
		"params":  []map[string]string{0: param},
	})
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return "", err
	}
	resp, err := http.Post(url, "application/json", strings.NewReader(string(data)))
	if err != nil {
		log.Fatalf("Post: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ReadAll: %v", err)
		return "", err
	}

	return string(body), nil
}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	initConf()
	if strings.Index(args.Url, "/iaas/list") == 0 {
		ListRunningVm(args, reply)
	}
	if strings.Index(args.Url, "/iaas/stop") == 0 {
		StopVm(args, reply)
	}
	if strings.Index(args.Url, "/iaas/start") == 0 {
		StartVm(args, reply)
	}
	if strings.Index(args.Url, "/iaas/download") == 0 {
		StartVm(args, reply)
	}
	if strings.Index(args.Url, "/iaas/downloadstatus") == 0 {
		StartVm(args, reply)
	}

	return nil
}

func (api) Plug(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Unplug(args interface{}, reply *bool) error {
	defer os.Exit(0)
	*reply = true
	return nil
}

func main() {

	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)
}
