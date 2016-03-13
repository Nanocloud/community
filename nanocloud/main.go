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
	"errors"
	"os"

	vmsConn "github.com/Nanocloud/community/nanocloud/connectors/vms"
	m "github.com/Nanocloud/community/nanocloud/middlewares"
	"github.com/Nanocloud/community/nanocloud/migration"
	appsModel "github.com/Nanocloud/community/nanocloud/models/apps"
	_ "github.com/Nanocloud/community/nanocloud/models/oauth"
	"github.com/Nanocloud/community/nanocloud/routes/apps"
	"github.com/Nanocloud/community/nanocloud/routes/front"
	"github.com/Nanocloud/community/nanocloud/routes/history"
	"github.com/Nanocloud/community/nanocloud/routes/logout"
	"github.com/Nanocloud/community/nanocloud/routes/machines"
	"github.com/Nanocloud/community/nanocloud/routes/me"
	"github.com/Nanocloud/community/nanocloud/routes/oauth"
	"github.com/Nanocloud/community/nanocloud/routes/tokens"
	"github.com/Nanocloud/community/nanocloud/routes/upload"
	"github.com/Nanocloud/community/nanocloud/routes/users"
	"github.com/Nanocloud/community/nanocloud/routes/version"
	"github.com/Nanocloud/community/nanocloud/utils"
	"github.com/Nanocloud/community/nanocloud/vms"
	_ "github.com/Nanocloud/community/nanocloud/vms/drivers/manual"
	_ "github.com/Nanocloud/community/nanocloud/vms/drivers/qemu"
	_ "github.com/Nanocloud/community/nanocloud/vms/drivers/vmwarefusion"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	logger "github.com/labstack/gommon/log"
)

func initVms() error {
	iaas := os.Getenv("IAAS")
	if len(iaas) == 0 {
		return errors.New("No iaas provided")
	}
	m := make(map[string]string, 0)

	switch iaas {

	case "vmwarefusion":
		m["PLAZA_LOCATION"] = os.Getenv("PLAZA_LOCATION")
		m["STORAGE_DIR"] = os.Getenv("STORAGE_DIR")

	default:
		m["servers"] = os.Getenv("EXECUTION_SERVERS")
		m["ad"] = os.Getenv("WIN_SERVER")
		m["sshport"] = os.Getenv("SSH_PORT")
		m["password"] = os.Getenv("WIN_PASSWORD")
		m["user"] = os.Getenv("WIN_USER")
	}

	vm, err := vms.Open(iaas, m)
	if err != nil {
		return err
	}
	vmsConn.SetVM(vm)
	return nil
}

func main() {
	err := migration.Migrate()
	if err != nil {
		log.Error(err)
		return
	}

	err = initVms()
	if err != nil {
		log.Error(err)
		return
	}

	e := echo.New()
	e.SetLogLevel(logger.DEBUG)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	/**
	 * LOGOUT
	 */
	e.Post("/api/logout", m.OAuth2(logout.Post))

	/**
	 * APPS
	 */
	e.Get("/api/apps", m.OAuth2(m.Admin(apps.ListApplications)))
	e.Delete("/api/apps/:app_id", m.OAuth2(m.Admin(apps.UnpublishApplication)))
	e.Get("/api/apps/me", m.OAuth2(apps.ListUserApps))
	e.Post("/api/apps", m.OAuth2(m.Admin(apps.PublishApplication)))
	e.Get("/api/apps/connections", m.OAuth2(apps.GetConnections))
	e.Patch("/api/apps/:app_id", m.OAuth2(m.Admin(apps.ChangeAppName)))

	go appsModel.CheckPublishedApps()

	/**
	 * HISTORY
	 */
	e.Get("/api/history", m.OAuth2(history.List))
	e.Post("/api/history", m.OAuth2(history.Add))

	/**
	 * USERS
	 */
	e.Patch("/api/users/:id", m.OAuth2(m.Admin(users.Update)))
	e.Get("/api/users", m.OAuth2(m.Admin(users.Get)))
	e.Post("/api/users", m.OAuth2(m.Admin(users.Post)))
	e.Delete("/api/users/:id", m.OAuth2(m.Admin(users.Delete)))
	e.Put("/api/users/:id", m.OAuth2(m.Admin(users.UpdatePassword)))
	e.Get("/api/users/:id", m.OAuth2(m.Admin(users.GetUser)))

	/**
	 * MACHINES
	 */
	e.Get("/api/machines", m.OAuth2(m.Admin(machines.Machines)))
	e.Get("/api/machines/:id", m.OAuth2(m.Admin(machines.GetMachine)))
	e.Patch("/api/machines/:id", m.OAuth2(m.Admin(machines.PatchMachine)))
	e.Post("/api/machines", m.OAuth2(m.Admin(machines.CreateMachine)))
	e.Delete("/api/machines/:id", m.OAuth2(m.Admin(machines.DeleteMachine)))

	/**
	 * FRONT
	 */
	e.Static("/canva/", front.StaticCanvaDirectory)
	e.Static("/", front.StaticDirectory)

	/**
	 * ME
	 */
	e.Get("/api/me", m.OAuth2(me.Get))

	/**
	 * VERSION
	 */
	e.Get("/api/version", version.Get)

	/**
	 * OAUTH
	 */
	e.Any("/oauth/*", oauth.Handler)

	/**
	 * TOKENS
	 */
	e.Get("/api/tokens", m.OAuth2(tokens.Get))
	e.Delete("/api/tokens/:id", m.OAuth2(tokens.Delete))

	/**
	 * UPLOAD
	 */
	e.Post("/upload", upload.Post)
	e.Get("/upload", upload.Get)

	addr := ":" + utils.Env("PORT", "8080")
	log.Info("Server running at ", addr)
	e.Run(addr)
}
