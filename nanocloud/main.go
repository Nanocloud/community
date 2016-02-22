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
	m "github.com/Nanocloud/community/nanocloud/middlewares"
	"github.com/Nanocloud/community/nanocloud/migration"
	appsModel "github.com/Nanocloud/community/nanocloud/models/apps"
	_ "github.com/Nanocloud/community/nanocloud/models/oauth"
	"github.com/Nanocloud/community/nanocloud/routes/apps"
	"github.com/Nanocloud/community/nanocloud/routes/front"
	"github.com/Nanocloud/community/nanocloud/routes/history"
	"github.com/Nanocloud/community/nanocloud/routes/iaas"
	"github.com/Nanocloud/community/nanocloud/routes/me"
	"github.com/Nanocloud/community/nanocloud/routes/oauth"
	"github.com/Nanocloud/community/nanocloud/routes/upload"
	"github.com/Nanocloud/community/nanocloud/routes/users"
	"github.com/Nanocloud/community/nanocloud/routes/version"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := migration.Migrate()
	if err != nil {
		log.Error(err)
		return
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	/**
	 * APPS
	 */
	e.Get("/api/apps", m.OAuth2(m.Admin(apps.ListApplications)))
	e.Delete("/api/apps/:app_id", m.OAuth2(m.Admin(apps.UnpublishApplication)))
	e.Get("/api/apps/me", m.OAuth2(apps.ListUserApps))
	e.Post("/api/apps", m.OAuth2(apps.PublishApplication))
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
	 * IAAS
	 */
	e.Get("/api/iaas", m.OAuth2(m.Admin(iaas.ListRunningVM)))
	e.Post("/api/iaas/:id/stop", m.OAuth2(m.Admin(iaas.StopVM)))
	e.Post("/api/iaas/:id/start", m.OAuth2(m.Admin(iaas.StartVM)))
	e.Post("/api/iaas/:id/download", m.OAuth2(m.Admin(iaas.DownloadVM)))

	/**
	 * FRONT
	 */
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
	 * UPLOAD
	 */
	e.Post("/upload", upload.Post)
	e.Get("/upload", upload.Get)

	addr := ":" + utils.Env("PORT", "8080")
	log.Info("Server running at ", addr)
	e.Run(addr)
}
