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
	appsModel "github.com/Nanocloud/community/nanocloud/models/apps"
	_ "github.com/Nanocloud/community/nanocloud/models/oauth"
	"github.com/Nanocloud/community/nanocloud/router"
	"github.com/Nanocloud/community/nanocloud/routes/apps"
	"github.com/Nanocloud/community/nanocloud/routes/front"
	"github.com/Nanocloud/community/nanocloud/routes/history"
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
	err := setupDb()
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
	router.Get("/apps", apps.ListApplications)
	router.Delete("/apps/:app_id", apps.UnpublishApplication)
	router.Get("/apps/me", apps.ListUserApps)
	router.Post("/apps", apps.PublishApplication)
	router.Get("/apps/connections", apps.GetConnections)
	router.Put("/apps/:app_id", middlewares.Admin, apps.ChangeAppName)

	go appsModel.CheckPublishedApps()

	/**
	 * HISTORY
	 */
	router.Get("/history", history.List)
	router.Post("/history", history.Add)

	/**
	 * USERS
	 */
	router.Post("/users/login", users.Login)
	router.Post("/users/:id/disable", middlewares.Admin, users.Disable)
	router.Get("/users", middlewares.Admin, users.Get)
	router.Post("/users", middlewares.Admin, users.Post)
	router.Delete("/users/:id", middlewares.Admin, users.Delete)
	router.Put("/users/:id", middlewares.Admin, users.UpdatePassword)
	router.Get("/users/:id", middlewares.Admin, users.GetUser)

	/**
	 * FRONT
	 */
	e.Static("/", front.StaticDirectory)

	/**
	 * ME
	 */
	e.Get("/api/me", me.Get)

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
