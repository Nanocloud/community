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
	"io/ioutil"
	"net/http"
	"strings"

	"fmt"
	"github.com/Nanocloud/community/nanocloud/models/apps"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

// ========================================================================================================================
// Procedure: createConnections
//
// Does:
// - Create all connections in DB for a particular user in order to use all applications
// ========================================================================================================================
func GetConnections(c *echo.Context) error {
	userList, err := users.FindUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": "Unable to retrieve users",
		})
	}
	user := c.Get("user").(*users.User)
	connections, err := apps.RetrieveConnections(user, userList)
	if err == apps.AppsListUnavailable {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": "Unable to retrieve applications list",
		})
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
	return c.JSON(http.StatusOK, hash{"data": response})
}

func ListApplications(c *echo.Context) error {
	user := c.Get("user").(*users.User)

	if !user.IsAdmin {
		applications, err := apps.GetUserApps(user.Id)
		if err == apps.GetAppsFailed {
			return c.JSON(http.StatusInternalServerError, hash{
				"error": [1]hash{
					hash{
						"detail": err.Error(),
					},
				},
			})
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
		return c.JSON(http.StatusOK, hash{"data": response})
	}

	applications, err := apps.GetAllApps()
	if err == apps.GetAppsFailed {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
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
	return c.JSON(http.StatusOK, hash{"data": response})
}

// Make an application unusable
func UnpublishApplication(c *echo.Context) error {
	appId := c.Param("app_id")
	if len(appId) < 1 {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "App id must be specified",
				},
			},
		})
	}

	err := apps.UnpublishApp(appId)
	if err == apps.UnpublishFailed {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, hash{
		"meta": hash{},
	})
}

func AddApplication(c *echo.Context) error {
	var p apps.ApplicationParams

	err := utils.ParseJSONBody(c, &p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "App infos are invalid",
				},
			},
		})
	}
	err = apps.AddApp(p)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err,
				},
			},
		})
	}
	return c.JSON(http.StatusOK, hash{
		"data": hash{
			"success": true,
		},
	})
}

func PublishApplication(c *echo.Context) error {
	var params struct {
		Data struct {
			Attributes struct {
				Path string `json:"path"`
			}
		}
	}

	err := utils.ParseJSONBody(c, &params)
	if err != nil {
		return nil
	}

	fmt.Printf("%+v\n", params)
	trimmedpath := strings.TrimSpace(params.Data.Attributes.Path)
	if trimmedpath == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "App path is empty",
				},
			},
		})
	}
	err = apps.PublishApp(trimmedpath)
	if err == apps.PublishFailed {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err,
				},
			},
		})
	}

	return c.JSON(http.StatusOK, hash{
		"data": hash{
			"success": true,
		},
	})
}

func ChangeAppName(c *echo.Context) error {
	appId := c.Param("app_id")
	if len(appId) < 1 {
		return c.JSON(http.StatusBadRequest, hash{
			"error": "App id must be specified",
		})
	}
	var Name struct {
		Data struct {
			Attributes struct {
				DisplayName string `json:"display-name"`
			}
		}
	}

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &Name)
	if err != nil {
		log.Errorf("Unable to parse body request: %s", err.Error())
		return err
	}
	if len(Name.Data.Attributes.DisplayName) < 1 {
		log.Errorf("No name provided")
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "No name provided",
				},
			},
		})
	}

	exists, err := apps.AppExists(appId)
	if err != nil {
		log.Errorf("Unable to check app existence: %s", err.Error())
		return err
	}

	if !exists {
		return c.JSON(http.StatusNotFound, hash{
			"error": [1]hash{
				hash{
					"detail": "App not found",
				},
			},
		})
	}

	err = apps.ChangeName(appId, Name.Data.Attributes.DisplayName)
	if err == apps.FailedNameChange {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
	}
	return c.JSON(http.StatusOK, hash{
		"data": hash{
			"success": true,
		},
	})
}
