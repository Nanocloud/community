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
	"github.com/Nanocloud/oauth"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func getMeHandler(c *echo.Context) error {
	user := oauth.GetUserOrFail(c.Response().Writer(), c.Request())
	if user == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, user)
}

func oauthHandler(w http.ResponseWriter, r *http.Request) {
	oauth.HandleRequest(w, r)
}

// get list of available front components
func getComponentsHandler(c *echo.Context) error {
	fis, err := ioutil.ReadDir(filepath.Join(env("FRONT_DIR", "front/"), "ts/components"))
	if err != nil {
		module.Log.Fatal("Unable to load the components folder. ", err)
		return c.Err()
	}
	var comps []string
	for _, f := range fis {
		comps = append(comps, f.Name())
	}
	return c.JSON(http.StatusOK, comps)
}

// get the version of the nanocloud application
func getVersionHandler(c *echo.Context) error {
	info := map[string]string{
		"version": appversion,
	}
	return c.JSON(http.StatusOK, info)
}
