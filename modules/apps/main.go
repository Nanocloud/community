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
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/Nanocloud/nano"
)

var module nano.Module

var conf struct {
	User                 string
	Server               string
	ExecutionServers     []string
	SSHPort              string
	RDPPort              string
	Password             string
	WindowsDomain        string
	XMLConfigurationFile string
}

type hash map[string]interface{}
type GetApplicationsListReply struct {
	Applications []Connection
}

func env(key, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		v = def
	}
	return v
}

// Make an application unusable
func unpublishApplication(req nano.Request) (*nano.Response, error) {
	appId := req.Params["app_id"]
	if len(appId) < 1 {
		return nano.JSONResponse(400, hash{
			"error": "App id must be specified",
		}), nil
	}

	err := unpublishApp(appId)
	if err != nil {
		return nil, err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func publishApplication(req nano.Request) (*nano.Response, error) {

	var params struct {
		Path string
	}

	err := json.Unmarshal(req.Body, &params)
	if err != nil {
		module.Log.Error("Umable to unmarshal application path: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": "path to app must be specified",
		}), err
	}

	trimmedpath := strings.TrimSpace(params.Path)
	if trimmedpath == "" {
		return nano.JSONResponse(400, hash{"error": "App path is empty"}), err
	}
	err = publishApp(trimmedpath)
	if err != nil {
		return nano.JSONResponse(500, hash{"error": err}), err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func main() {
	module = nano.RegisterModule("apps")

	conf.User = env("USER", "Administrator")
	conf.SSHPort = env("SSH_PORT", "22")
	conf.RDPPort = env("RDP_PORT", "3389")
	conf.Server = env("SERVER", "62.210.56.76")
	conf.Password = env("PASSWORD", "ItsPass1942+")
	conf.XMLConfigurationFile = env("XML_CONFIGURATION_FILE", "conf.xml")
	conf.WindowsDomain = env("WINDOWS_DOMAIN", "intra.localdomain.com")
	conf.ExecutionServers = strings.Split(env("EXECUTION_SERVERS", "62.210.56.76"), ",")

	module.Get("/apps", listApplications)
	module.Delete("/apps/:app_id", unpublishApplication)
	module.Get("/apps/me", listApplicationsForSamAccount)
	module.Post("/apps", publishApplication)
	err := errors.New("")
	for err != nil {
		err = createConnections()
		module.Log.Error(err)
		time.Sleep(time.Second * 3)
	}

	module.Listen()
}
