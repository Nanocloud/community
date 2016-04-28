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

package histories

import (
	"github.com/Nanocloud/community/nanocloud/models/histories"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"net/http"
)

type hash map[string]interface{}

// Get a list of all the log entries of the database
func List(c *echo.Context) error {

	histories, err := histories.FindAll()
	if err != nil {
		return err
	}
	return utils.JSON(c, http.StatusOK, histories)
}

// Add a new log entry to the database
func Add(c *echo.Context) error {
	history := histories.History{}

	err := utils.ParseJSONBody(c, &history)
	if err != nil {
		return err
	}

	if history.UserId == "" || history.ConnectionId == "" || history.StartDate == "" || history.EndDate == "" {
		log.Error("Missing one or several parameters to create entry")
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "Missing parameters",
				},
			},
		})
	}

	err = utils.ParseJSONBody(c, &history)
	newHistory, err := histories.CreateHistory(
		history.UserId,
		history.ConnectionId,
		history.StartDate,
		history.EndDate,
	)

	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusCreated, newHistory)
}
