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

package middlewares

import (
	"errors"

	"github.com/Nanocloud/community/nanocloud/oauth2"
	"gopkg.in/labstack/echo.v1"
)

func oAuth2(c *echo.Context, handler echo.HandlerFunc) error {
	r := c.Request()
	w := c.Response()

	user, err := oauth2.GetUser(w, r)
	if err != nil {
		b, fail := err.ToJSON()
		if fail != nil {
			return fail
		}

		w.WriteHeader(err.HTTPStatusCode)
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return nil
	}

	if user != nil {
		c.Set("user", user)
		return handler(c)
	}
	return errors.New("unable to authenticate the user")
}

func OAuth2(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return oAuth2(c, handler)
	}
}
