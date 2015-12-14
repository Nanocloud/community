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
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"regexp"

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

// Plugin structure
type api struct{}

type VmName struct {
	Name string
}

// Structure used for exchanges with the core, faking a request/responsewriter
type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
	Method   string
	HeadVals map[string]string
	Status   int
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ListRunningVm(args PlugRequest, reply *PlugRequest, name string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	reply.Status = 200
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
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.GetList",
		nil,
	)
	if err != nil {
		log.Error("RPC Call to Iaas API failed: ", err)
		reply.Status = 500
		return
	}

	err = json.Unmarshal([]byte(jsonResponse), &response)
	if err != nil {
		log.Error("Failed to Unmarshal response from Iaas API: ", err)
		reply.Status = 500
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

	jsonOutput, err := json.Marshal(vmList)
	if err != nil {
		log.Error("Failed to Marshal Vm list: ", err)
		reply.Status = 500
	}
	reply.Body = string(jsonOutput)
}

func DownloadVm(args PlugRequest, reply *PlugRequest, name string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	reply.Status = 200
	var params = map[string]string{
		"vmname": name,
	}

	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.Download",
		params,
	)
	if err != nil {
		log.Error("RPC Call to Iaas API failed: ", err)
		reply.Status = 500
		return
	}
	reply.Body = string(jsonResponse)

}

// NOT USEFUL RIGHT NOW
/*
func DownloadStatus(args PlugRequest, reply *PlugRequest, name string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	reply.Status = 200

	var params = map[string]string{
		"Name": name,
	}
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.GetStatus",
		params,
	)
	if err != nil {
		log.Error("RPC Call to Iaas API failed: ", err)
		reply.Status = 500
		return
	}
	reply.Body = string(jsonResponse)

}*/

func StartVm(args PlugRequest, reply *PlugRequest, name string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	reply.Status = 200

	var params = map[string]string{
		"name": name,
	}
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.Start",
		params,
	)
	if err != nil {
		log.Error("RPC Call to Iaas API failed: ", err)
		reply.Status = 500
		return
	}
	reply.Body = string(jsonResponse)

}

func StopVm(args PlugRequest, reply *PlugRequest, name string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	reply.Status = 200

	var params = map[string]string{
		"Name": name,
	}
	jsonResponse, err := jsonRpcRequest(
		fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		"Iaas.Stop",
		params,
	)
	if err != nil {
		log.Error("RPC Call to Iaas API failed: ", err)
		reply.Status = 500
		return
	}
	reply.Body = string(jsonResponse)
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

var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string)
}{
	{`^\/api\/iaas\/{0,1}$`, "GET", ListRunningVm},
	{`^\/api\/iaas\/(?P<id>[^\/]+)\/stop\/{0,1}$`, "POST", StopVm},
	{`^\/api\/iaas\/(?P<id>[^\/]+)\/start\/{0,1}$`, "POST", StartVm},
	{`^\/api\/iaas\/(?P<id>[^\/]+)\/download\/{0,1}$`, "POST", DownloadVm},
	//	{`^\/api\/iaas\/(?P<id>[^\/]+)\/status\/{0,1}$`, "POST", DownloadStatus},
}

// Will receive all http requests starting by /api/history from the core and chose the correct handler function
func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
			} else {
				val.f(args, reply, "")
			}
		}
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
	initConf()

	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)
}
