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

package utils

import (
	"io/ioutil"

	"github.com/Nanocloud/community/nanocloud/errors"
	"github.com/labstack/echo"
	"github.com/manyminds/api2go/jsonapi"
)

func JSON(c *echo.Context, code int, i interface{}) error {
	b, err := jsonapi.Marshal(i)

	if err != nil {
		return err
	}

	r := c.Response()

	r.Header().Set("Content-Type", "application/vnd.api+json")
	r.WriteHeader(code)
	r.Write(b)
	return nil
}

func ParseJSONBody(c *echo.Context, dest interface{}) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return errors.InvalidRequest
	}

	err = jsonapi.Unmarshal(body, dest)
	if err != nil {
		return errors.InvalidRequest
	}
	return nil
}
