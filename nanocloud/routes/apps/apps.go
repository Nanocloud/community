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

package apps

import (
	"encoding/json"
	"github.com/Nanocloud/community/nanocloud/models/apps"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/router"
	log "github.com/Sirupsen/logrus"
	"strings"
)

type hash map[string]interface{}

// ========================================================================================================================
// Procedure: createConnections
//
// Does:
// - Create all connections in DB for a particular user in order to use all applications
// ========================================================================================================================
func GetConnections(req router.Request) (*router.Response, error) {
	userList, err := users.FindUsers()
	if err != nil {
		return router.JSONResponse(500, hash{
			"error": "Unable to retrieve users",
		}), nil
	}
	connections, err := apps.RetrieveConnections(userList)
	if err == apps.AppsListUnavailable {
		return router.JSONResponse(500, hash{
			"error": "Unable to retrieve applications list",
		}), nil
	}
	return router.JSONResponse(200, connections), nil
}

func ListApplications(req router.Request) (*router.Response, error) {
	applications, err := apps.GetAllApps()
	if err == apps.GetAppsFailed {
		return router.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil

	}
	return router.JSONResponse(200, applications), nil
}

// ========================================================================================================================
// Procedure: listApplicationsForSamAccount
//
// Does:
// - Return list of applications available for a particular SAM account
// ========================================================================================================================
func ListApplicationsForSamAccount(req router.Request) (*router.Response, error) {
	applications, err := apps.GetMyApps()
	if err == apps.GetAppsFailed {
		return router.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	//TODO ONLY RETURN AUTHORIZED APPS FROM THE USER'S GROUP
	return router.JSONResponse(200, applications), nil
}

// Make an application unusable
func UnpublishApplication(req router.Request) (*router.Response, error) {
	appId := req.Params["app_id"]
	if len(appId) < 1 {
		return router.JSONResponse(400, hash{
			"error": "App id must be specified",
		}), nil
	}

	err := apps.UnpublishApp(appId)
	if err == apps.UnpublishFailed {
		return router.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	return router.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func PublishApplication(req router.Request) (*router.Response, error) {
	var params struct {
		Path string
	}

	err := json.Unmarshal(req.Body, &params)
	if err != nil {
		log.Error("Umable to unmarshal application path: " + err.Error())
		return router.JSONResponse(400, hash{
			"error": "path to app must be specified",
		}), err
	}

	trimmedpath := strings.TrimSpace(params.Path)
	if trimmedpath == "" {
		return router.JSONResponse(400, hash{"error": "App path is empty"}), err
	}
	err = apps.PublishApp(trimmedpath)
	if err == apps.PublishFailed {
		return router.JSONResponse(500, hash{"error": err}), err
	}

	return router.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func ChangeAppName(req router.Request) (*router.Response, error) {
	appId := req.Params["app_id"]
	if len(appId) < 1 {
		return router.JSONResponse(400, hash{
			"error": "App id must be specified",
		}), nil
	}
	var Name struct {
		DisplayName string
	}

	err := json.Unmarshal(req.Body, &Name)
	if err != nil {
		log.Errorf("Unable to parse body request: %s", err.Error())
		return nil, err
	}
	if len(Name.DisplayName) < 1 {
		log.Errorf("No name provided")
		return router.JSONResponse(400, hash{
			"error": "No name provided",
		}), nil
	}

	exists, err := apps.AppExists(appId)
	if err != nil {
		log.Errorf("Unable to check app existence: %s", err.Error())
		return nil, err
	}

	if !exists {
		return router.JSONResponse(404, hash{
			"error": "App not found",
		}), nil
	}

	err = apps.ChangeName(appId, Name.DisplayName)
	if err == apps.FailedNameChange {
		return router.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}
	return router.JSONResponse(200, hash{
		"success": true,
	}), nil
}
