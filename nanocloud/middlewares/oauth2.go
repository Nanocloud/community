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
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/oauth2"
	"github.com/Nanocloud/community/nanocloud/router"
)

func OAuth2(req *router.Request) (*router.Response, error) {
	r := req.Request()
	w := req.Response()

	user, err := oauth2.GetUser(w, r)
	if err != nil {
		b, fail := err.ToJSON()
		if fail != nil {
			return nil, fail
		}

		res := router.Response{
			ContentType: "application/json",
			Body:        b,
			StatusCode:  err.HTTPStatusCode,
		}

		return &res, nil
	}

	req.User = user.(*users.User)
	return nil, nil
}
