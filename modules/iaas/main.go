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

	"github.com/Nanocloud/community/modules/iaas/lib/iaas"
	"github.com/Nanocloud/nano"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
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

func (h *handler) ListRunningVM(c *echo.Context) error {
	response, err := h.iaasCon.GetList()
	if err != nil {
		log.Error("Unable to retrieve VM states list")
		return c.JSON(http.StatusInternalServerError, hash{
			"error": "Unable te retrieve states of VMs: " + err.Error(),
		})
	}

	vmList := h.iaasCon.CheckVMStates(response)
	return c.JSON(http.StatusOK, vmList)
}

func (h *handler) DownloadVM(c *echo.Context) error {
	vmname := c.Param("id")

	if vmname == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": "No VM ID provided",
		})
	}

	go h.iaasCon.Download(vmname)
	return c.JSON(http.StatusOK, hash{
		"success": true,
	})
}

func (h *handler) StartVM(c *echo.Context) error {
	name := c.Param("id")

	if name == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": "No VM name provided",
		})
	}

	err := h.iaasCon.Start(name)
	if err != nil {
		log.Error("Error while starting VM")
		return c.JSON(http.StatusInternalServerError, hash{
			"error": "Unable to start the specified VM",
		})
	}

	return c.JSON(http.StatusOK, hash{
		"success": true,
	})
}

func (h *handler) StopVM(c *echo.Context) error {
	name := c.Param("id")

	if name == "" {
		return c.JSON(http.StatusBadRequest, hash{
			"error": "No VM name provided",
		})
	}

	err := h.iaasCon.Stop(name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, hash{
			"error": "Unable to stop the specified VM",
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

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	h := handler{
		iaasCon: iaas.New(conf.Server, conf.Password, conf.User, conf.SSHPort, conf.instDir, conf.artURL),
	}

	e.Get("/api/iaas", h.ListRunningVM)
	e.Post("/api/iaas/:id/stop", h.StopVM)
	e.Post("/api/iaas/:id/start", h.StartVM)
	e.Post("/api/iaas/:id/download", h.DownloadVM)

	e.Run(":8080")
}
