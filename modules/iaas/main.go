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

	"github.com/Nanocloud/community/modules/iaas/lib/iaas"
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

type hash map[string]interface{}

type handler struct {
	iaasCon *iaas.Iaas
}

var conf Configuration

func AdminOnly(req nano.Request) (*nano.Response, error) {
	if req.User != nil && !req.User.IsAdmin {
		return nano.JSONResponse(403, hash{
			"error": "forbidden",
		}), nil
	}
	return nil, nil
}

func (h *handler) ListRunningVM(req nano.Request) (*nano.Response, error) {
	response, err := h.iaasCon.GetList()
	if err != nil {
		module.Log.Error("Unable to retrieve VM states list")
		return nano.JSONResponse(500, hash{
			"error": "Unable te retrieve states of VMs: " + err.Error(),
		}), err
	}

	vmList := h.iaasCon.CheckVMStates(response)
	return nano.JSONResponse(200, vmList), nil
}

func (h *handler) DownloadVM(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"vmname": req.Params["id"],
	}

	if params["vmname"] == "" {
		return nano.JSONResponse(400, hash{
			"error": "No VM ID provided",
		}), nil
	}

	go h.iaasCon.Download(params["vmname"])
	return nano.JSONResponse(202, hash{
		"success": true,
	}), nil
}

func (h *handler) StartVM(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"name": req.Params["id"],
	}

	if params["name"] == "" {
		return nano.JSONResponse(400, hash{
			"error": "No VM name provided",
		}), nil
	}

	err := h.iaasCon.Start(params["name"])
	if err != nil {
		module.Log.Error("Error while starting VM")
		return nano.JSONResponse(500, hash{
			"error": "Unable to start the specified VM",
		}), err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func (h *handler) StopVM(req nano.Request) (*nano.Response, error) {
	var params = map[string]string{
		"Name": req.Params["id"],
	}

	if params["Name"] == "" {
		return nano.JSONResponse(400, hash{
			"error": "No VM name provided",
		}), nil
	}

	err := h.iaasCon.Stop(params["Name"])
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": "Unable to stop the specified VM",
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

	conf.Server = env("SERVER", "127.0.0.1")
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

	h := handler{
		iaasCon: iaas.New(conf.Server, conf.Password, conf.User, conf.SSHPort, conf.instDir, conf.artURL),
	}

	module.Get("/iaas", AdminOnly, h.ListRunningVM)
	module.Post("/iaas/:id/stop", AdminOnly, h.StopVM)
	module.Post("/iaas/:id/start", AdminOnly, h.StartVM)
	module.Post("/iaas/:id/download", AdminOnly, h.DownloadVM)

	module.Listen()
}
