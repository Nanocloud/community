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

package tokens

import (
	"fmt"
	"net/http"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

func Get(c *echo.Context) error {
	user := c.Get("user").(*users.User)

	res, err := db.Query(
		"SELECT id, created_at FROM oauth_access_tokens WHERE user_id = $1::varchar",
		user.Id,
	)

	if err != nil {
		return err
	}

	defer res.Close()

	r := make([]hash, 0)
	for res.Next() {
		var id, createdAt string
		err := res.Scan(&id, &createdAt)
		if err != nil {
			continue
		}

		r = append(r, hash{
			"id":   id,
			"type": "token",
			"attributes": hash{
				"created-at": createdAt,
			},
		})
	}

	return c.JSON(http.StatusOK, r)
}

func Delete(c *echo.Context) error {
	tokenId := c.Param("id")

	if len(tokenId) == 4 {
		return c.JSON(http.StatusBadRequest, hash{
			"error": "Invalid token id",
		})
	}

	user := c.Get("user").(*users.User)

	fmt.Println(user.Id)
	fmt.Println(tokenId)
	res, err := db.Exec(
		`DELETE FROM oauth_access_tokens
		WHERE user_id = $1::varchar
		AND id = $2::varchar`,
		user.Id,
		tokenId,
	)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected < 1 {
		return c.JSON(http.StatusNotFound, hash{
			"error": "No such token",
		})
	}

	return c.JSON(http.StatusOK, hash{})
}
