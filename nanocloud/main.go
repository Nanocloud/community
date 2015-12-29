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
	"github.com/Nanocloud/nano"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

var module nano.Module

func env(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		v = def
	}
	return v
}

func main() {
	setupDb()

	module = nano.RegisterModule("router")

	handler := httpHandler{
		URLPrefix: "/api",
		Module:    module,
	}

	go module.Listen()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	frontDir := env("FRONT_DIR", "front/")
	e.Static("/", frontDir)
	e.Get("/api/me", getMeHandler)
	e.Get("/api/version", getVersionHandler)
	e.Any("/api/*", handler.ServeHTTP)
	e.Any("/oauth/*", oauthHandler)
	e.Post("/upload", uploadHandler)
	e.Get("/upload", checkUploadHandler)

	addr := ":" + env("PORT", "8080")
	module.Log.Info("Server running at ", addr)
	e.Run(addr)
}
