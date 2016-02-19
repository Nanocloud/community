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
	"strings"

	"github.com/Nanocloud/community/nanocloud/models/apps"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/router"
	log "github.com/Sirupsen/logrus"
)

type hash map[string]interface{}

// ========================================================================================================================
// Procedure: createConnections
//
// Does:
// - Create all connections in DB for a particular user in order to use all applications
// ========================================================================================================================
func GetConnections(req *router.Request) (*router.Response, error) {
	userList, err := users.FindUsers()
	if err != nil {
		return router.JSONResponse(500, hash{
			"error": "Unable to retrieve users",
		}), nil
	}
	connections, err := apps.RetrieveConnections(req.User, userList)
	if err == apps.AppsListUnavailable {
		return router.JSONResponse(500, hash{
			"error": "Unable to retrieve applications list",
		}), nil
	}

	var response = make([]hash, len(connections))
	for i, val := range connections {
		res := hash{
			"id":         i,
			"type":       "application",
			"attributes": val,
		}
		response[i] = res
	}
	return router.JSONResponse(200, hash{"data": response}), nil
}

func ListApplications(req *router.Request) (*router.Response, error) {
	applications, err := apps.GetAllApps()
	if err == apps.GetAppsFailed {
		return router.JSONResponse(500, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		}), nil
	}

	var response = make([]hash, len(applications))
	for i, val := range applications {
		res := hash{
			"id":         val.Id,
			"type":       "application",
			"attributes": val,
		}
		response[i] = res
	}
	return router.JSONResponse(200, hash{"data": response}), nil
}

// ========================================================================================================================
// Procedure: ListUserApps
//
// Does:
// - Return list of applications available for the current user
// ========================================================================================================================
func ListUserApps(req *router.Request) (*router.Response, error) {
	applications, err := apps.GetUserApps(req.User.Id)
	if err == apps.GetAppsFailed {
		return router.JSONResponse(500, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		}), nil
	}

	var response = make([]hash, len(applications))
	for i, val := range applications {
		res := hash{
			"id":         val.Id,
			"type":       "application",
			"attributes": val,
		}
		response[i] = res
	}
	return router.JSONResponse(200, hash{"data": response}), nil
}

// Make an application unusable
func UnpublishApplication(req *router.Request) (*router.Response, error) {
	appId := req.Params["app_id"]
	if len(appId) < 1 {
		return router.JSONResponse(400, hash{
			"error": [1]hash{
				hash{
					"detail": "App id must be specified",
				},
			},
		}), nil
	}

	err := apps.UnpublishApp(appId)
	if err == apps.UnpublishFailed {
		return router.JSONResponse(500, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		}), nil
	}

	return router.JSONResponse(200, hash{
		"data": hash{
			"success": true,
		},
	}), nil
}

func PublishApplication(req *router.Request) (*router.Response, error) {
	var params struct {
		Data struct {
			Attributes struct {
				Path string `json:"path"`
			}
		}
	}

	err := json.Unmarshal(req.Body, &params)
	if err != nil {
		log.Error("Umable to unmarshal application path: " + err.Error())
		return router.JSONResponse(400, hash{
			"error": [1]hash{
				hash{
					"detail": "path to app must be specified",
				},
			},
		}), err
	}

	trimmedpath := strings.TrimSpace(params.Data.Attributes.Path)
	if trimmedpath == "" {
		return router.JSONResponse(400, hash{
			"error": [1]hash{
				hash{
					"detail": "App path is empty",
				},
			},
		}), err
	}
	err = apps.PublishApp(trimmedpath)
	if err == apps.PublishFailed {
		return router.JSONResponse(500, hash{
			"error": [1]hash{
				hash{
					"detail": err,
				},
			},
		}), err
	}

	return router.JSONResponse(200, hash{
		"data": hash{
			"success": true,
		},
	}), nil
}

func ChangeAppName(req *router.Request) (*router.Response, error) {
	appId := req.Params["app_id"]
	if len(appId) < 1 {
		return router.JSONResponse(400, hash{
			"error": "App id must be specified",
		}), nil
	}
	var Name struct {
		Data struct {
			DisplayName string `json:"display_name"`
		}
	}

	err := json.Unmarshal(req.Body, &Name)
	if err != nil {
		log.Errorf("Unable to parse body request: %s", err.Error())
		return nil, err
	}
	if len(Name.Data.DisplayName) < 1 {
		log.Errorf("No name provided")
		return router.JSONResponse(400, hash{
			"error": [1]hash{
				hash{
					"detail": "No name provided",
				},
			},
		}), nil
	}

	exists, err := apps.AppExists(appId)
	if err != nil {
		log.Errorf("Unable to check app existence: %s", err.Error())
		return nil, err
	}

	if !exists {
		return router.JSONResponse(404, hash{
			"error": [1]hash{
				hash{
					"detail": "App not found",
				},
			},
		}), nil
	}

	err = apps.ChangeName(appId, Name.Data.DisplayName)
	if err == apps.FailedNameChange {
		return router.JSONResponse(500, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		}), nil
	}
	return router.JSONResponse(200, hash{
		"data": hash{
			"success": true,
		},
	}), nil
}
