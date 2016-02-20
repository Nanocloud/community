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
	"github.com/Nanocloud/community/nanocloud/middlewares"
	"github.com/Nanocloud/community/nanocloud/migration"
	appsModel "github.com/Nanocloud/community/nanocloud/models/apps"
	_ "github.com/Nanocloud/community/nanocloud/models/oauth"
	"github.com/Nanocloud/community/nanocloud/router"
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
	e.Any("/api/*", router.ServeHTTP)

	/**
	 * APPS
	 */
	router.Get("/apps", middlewares.OAuth2, apps.ListApplications)
	router.Delete("/apps/:app_id", middlewares.OAuth2, apps.UnpublishApplication)
	router.Get("/apps/me", middlewares.OAuth2, apps.ListUserApps)
	router.Post("/apps", middlewares.OAuth2, apps.PublishApplication)
	router.Get("/apps/connections", middlewares.OAuth2, apps.GetConnections)
	router.Patch("/apps/:app_id", middlewares.OAuth2, middlewares.Admin, apps.ChangeAppName)

	go appsModel.CheckPublishedApps()

	/**
	 * HISTORY
	 */
	router.Get("/history", middlewares.OAuth2, history.List)
	router.Post("/history", middlewares.OAuth2, history.Add)

	/**
	 * USERS
	 */
	router.Patch("/users/:id", middlewares.OAuth2, middlewares.Admin, users.Update)
	router.Get("/users", middlewares.OAuth2, middlewares.Admin, users.Get)
	router.Post("/users", middlewares.OAuth2, middlewares.Admin, users.Post)
	router.Delete("/users/:id", middlewares.OAuth2, middlewares.Admin, users.Delete)
	router.Put("/users/:id", middlewares.OAuth2, middlewares.Admin, users.UpdatePassword)
	router.Get("/users/:id", middlewares.OAuth2, middlewares.Admin, users.GetUser)

	/**
	 * IAAS
	 */
	router.Get("/iaas", middlewares.OAuth2, middlewares.Admin, iaas.ListRunningVM)
	router.Post("/iaas/:id/stop", middlewares.OAuth2, middlewares.Admin, iaas.StopVM)
	router.Post("/iaas/:id/start", middlewares.OAuth2, middlewares.Admin, iaas.StartVM)
	router.Post("/iaas/:id/download", middlewares.OAuth2, middlewares.Admin, iaas.DownloadVM)

	/**
	 * FRONT
	 */
	e.Static("/", front.StaticDirectory)

	/**
	 * ME
	 */
	router.Get("/me", middlewares.OAuth2, me.Get)

	/**
	 * VERSION
	 */
	router.Get("/version", version.Get)

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
