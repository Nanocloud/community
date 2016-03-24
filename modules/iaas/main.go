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
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type Configuration struct {
	artURL     string
	instDir    string
	Server     string
	User       string
	SSHPort    string
	Password   string
	windowsURL string
}

type hash map[string]interface{}

var conf Configuration

func ListRunningVM(c *echo.Context) error {
	response, err := GetList()
	if err != nil {
		log.Error("Unable to retrieve VM states list")
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": "Unable te retrieve states of VMs: " + err.Error(),
				},
			},
		})
	}

	vmList := CheckVMStates(response)
	var res = make([]hash, len(vmList))
	for i, val := range vmList {
		r := hash{
			"id":         val.Name,
			"type":       "vm",
			"attributes": val,
		}
		res[i] = r
	}
	return c.JSON(http.StatusOK, hash{"data": res})
}

func DownloadVM(c *echo.Context) error {
	vmname := c.Param("id")

	if vmname == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "No VM ID provided",
				},
			},
		})
	}

	go Download(vmname)
	return c.JSON(http.StatusOK, hash{
		"success": true,
	})
}

func StartVM(c *echo.Context) error {
	name := c.Param("id")

	if name == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "No VM name provided",
				},
			},
		})
	}

	err := Start(name)
	if err != nil {
		log.Error("Error while starting VM")
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
	}

	return c.JSON(http.StatusOK, hash{
		"success": true,
	})
}

func StopVM(c *echo.Context) error {
	name := c.Param("id")

	if name == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "No VM name provided",
				},
			},
		})
	}

	err := Stop(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": "Unable to stop the specified VM",
				},
			},
		})
	}

	return c.JSON(http.StatusOK, hash{
		"success": true,
	})
}

func CreateVM(c *echo.Context) error {

	err := Create()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": [1]hash{
				hash{
					"detail": "Unable to create the specified VM",
				},
			},
		})
	}

	return c.JSON(http.StatusOK, hash{
		"success": true,
	})
}

func env(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func main() {
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

	conf.windowsURL = os.Getenv("WINDOWS_URL")
	if len(conf.windowsURL) == 0 {
		conf.windowsURL = "http://care.dlservice.microsoft.com/dl/download/6/2/A/62A76ABB-9990-4EFC-A4FE-C7D698DAEB96/9600.17050.WINBLUE_REFRESH.140317-1640_X64FRE_SERVER_EVAL_EN-US-IR3_SSS_X64FREE_EN-US_DV9.ISO"
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.Get("/api/vms", ListRunningVM)
	e.Post("/api/vms/:id/stop", StopVM)
	e.Post("/api/vms/:id/start", StartVM)
	e.Post("/api/vms/:id/download", CreateVM)

	e.Run(":8080")
}
