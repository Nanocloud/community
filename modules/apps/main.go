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

	"github.com/Nanocloud/community/modules/apps/lib/apps"
	"github.com/Nanocloud/nano"
	_ "github.com/lib/pq"
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
	DatabaseURI          string
	Protocol             string
}

type handler struct {
	appsCon *apps.Apps
}

type hash map[string]interface{}

func env(key, def string) string {
	v := os.Getenv(key)
	if len(v) == 0 {
		v = def
	}
	return v
}

func getUsers() ([]nano.User, error) {
	res, err := module.Request("GET", "/users", "", nil, nil)
	if err != nil {
		module.Log.Error("GET on /users failed: ", err)
		return nil, err
	}

	if res.StatusCode != 200 {
		module.Log.Error("Status code of GET on /users isn't 200")
		return nil, errors.New("invalid status code")
	}

	var m []nano.User
	err = json.Unmarshal(res.Body, &m)

	if err != nil {
		module.Log.Error("Unmarshal of user list failed: ", err)
		return nil, err
	}

	return m, nil
}

// ========================================================================================================================
// Procedure: createConnections
//
// Does:
// - Create all connections in DB for a particular user in order to use all applications
// ========================================================================================================================
func (h *handler) getConnections(req nano.Request) (*nano.Response, error) {
	users, err := getUsers()
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": "Unable to retrieve users",
		}), nil
	}
	connections, err := h.appsCon.RetrieveConnections(users)
	if err == apps.AppsListUnavailable {
		return nano.JSONResponse(500, hash{
			"error": "Unable to retrieve applications list",
		}), nil
	}
	return nano.JSONResponse(200, connections), nil
}

func (h *handler) listApplications(req nano.Request) (*nano.Response, error) {

	applications, err := h.appsCon.GetAllApps()
	if err == apps.GetAppsFailed {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil

	}
	return nano.JSONResponse(200, applications), nil
}

// ========================================================================================================================
// Procedure: listApplicationsForSamAccount
//
// Does:
// - Return list of applications available for a particular SAM account
// ========================================================================================================================
func (h *handler) listApplicationsForSamAccount(req nano.Request) (*nano.Response, error) {
	applications, err := h.appsCon.GetMyApps()
	if err == apps.GetAppsFailed {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	//TODO ONLY RETURN AUTHORIZED APPS FROM THE USER'S GROUP
	return nano.JSONResponse(200, applications), nil
}

// Make an application unusable
func (h *handler) unpublishApplication(req nano.Request) (*nano.Response, error) {
	appId := req.Params["app_id"]
	if len(appId) < 1 {
		return nano.JSONResponse(400, hash{
			"error": "App id must be specified",
		}), nil
	}

	err := h.appsCon.UnpublishApp(appId)
	if err == apps.UnpublishFailed {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func (h *handler) publishApplication(req nano.Request) (*nano.Response, error) {

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
	err = h.appsCon.PublishApp(trimmedpath)
	if err == apps.PublishFailed {
		return nano.JSONResponse(500, hash{"error": err}), err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func main() {
	module = nano.RegisterModule("apps")

	conf.User = env("USER", "Administrator")
	conf.Protocol = env("PROTOCOL", "rdp")
	conf.SSHPort = env("SSH_PORT", "22")
	conf.RDPPort = env("RDP_PORT", "3389")
	conf.Server = env("SERVER", "62.210.56.76")
	conf.Password = env("PASSWORD", "ItsPass1942+")
	conf.ExecutionServers = strings.Split(env("EXECUTION_SERVERS", "62.210.56.76"), ",")
	conf.DatabaseURI = env("DATABASE_URI", "postgres://localhost/nanocloud?sslmode=disable")

	h := handler{
		appsCon: apps.New(
			conf.User,
			conf.Server,
			conf.SSHPort,
			conf.RDPPort,
			conf.Protocol,
			conf.Password,
			conf.DatabaseURI,
			conf.ExecutionServers,
		),
	}

	module.Get("/apps", h.listApplications)
	module.Delete("/apps/:app_id", h.unpublishApplication)
	module.Get("/apps/me", h.listApplicationsForSamAccount)
	module.Post("/apps", h.publishApplication)
	module.Get("/apps/connections", h.getConnections)
	module.Listen()
}
