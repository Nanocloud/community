/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2016 Nanocloud Software
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

package exec

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

type bodyRequest struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Domain   string   `json:"domain"`
	Command  []string `json:"command"`
	Stdin    string   `json:"stdin"`
	AppMode  bool     `json:"app-mode"`
}

func Route(c *echo.Context) error {
	b, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	body := bodyRequest{}
	err = json.Unmarshal(b, &body)
	if err != nil {
		return err
	}

	if body.AppMode {
		pid := uint32(0)
		pid, err = launchApp(body.Command)
		if err != nil {
			return err
		}

		m := make(map[string]uint32)
		m["pid"] = pid
		return c.JSON(http.StatusOK, m)
	}

	cmd := runCommand(body.Username, body.Domain, body.Password, body.Command)
	if body.Stdin != "" {
		cmd.Stdin = strings.NewReader(body.Stdin)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()

	res := makeResponse(stdout, stderr, cmd)
	return c.JSON(http.StatusOK, res)
}
