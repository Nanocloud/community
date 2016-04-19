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

package history

import (
	"net/http"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

// Log entries are stored in this structure
type HistoryInfo struct {
	UserId       string `json:"user-id"`
	ConnectionId string `json:"connection-id"`
	StartDate    string `json:"start-date"`
	EndDate      string `json:"end-date"`
}

// Get a list of all the log entries of the database
func List(c *echo.Context) error {
	var histories []HistoryInfo
	rows, err := db.Query(
		`SELECT userid, connectionid,
		startdate, enddate
		FROM histories`,
	)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		history := HistoryInfo{}

		rows.Scan(
			&history.UserId,
			&history.ConnectionId,
			&history.StartDate,
			&history.EndDate,
		)
		histories = append(histories, history)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	if len(histories) == 0 {
		histories = []HistoryInfo{}
	}

	var response = make([]hash, len(histories))
	for i, val := range histories {
		res := hash{
			"id":         i,
			"type":       "history",
			"attributes": val,
		}
		response[i] = res
	}
	return c.JSON(http.StatusOK, hash{"data": response})
}

// Add a new log entry to the database
func Add(c *echo.Context) error {
	var attr hash

	err := utils.ParseJSONBody(c, &attr)
	if err != nil {
		return err
	}

	data, ok := attr["data"].(map[string]interface{})
	if ok == false {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "data is missing",
				},
			},
		})
	}

	attributes, ok := data["attributes"].(map[string]interface{})
	if ok == false {
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "attributes is missing",
				},
			},
		})
	}

	user_id, ok := attributes["user_id"].(string)
	connection_id, ok := attributes["connection_id"].(string)
	start_date, ok := attributes["start_date"].(string)
	end_date, ok := attributes["end_date"].(string)
	if user_id == "" || connection_id == "" || start_date == "" || end_date == "" {
		log.Error("Missing one or several parameters to create entry")
		return c.JSON(http.StatusBadRequest, hash{
			"error": [1]hash{
				hash{
					"detail": "Missing parameters",
				},
			},
		})
	}

	rows, err := db.Query(
		`INSERT INTO histories
		(userid, connectionid, startdate, enddate)
		VALUES(	$1::varchar, $2::varchar, $3::varchar, $4::varchar)
		`, user_id, connection_id, start_date, end_date)
	if err != nil {
		return err
	}

	rows.Close()

	return c.JSON(http.StatusCreated, hash{
		"data": hash{
			"attributes": hash{
				"user_id":       user_id,
				"connection_id": connection_id,
				"start_date":    start_date,
				"end_date":      end_date,
			},
			"type": "history",
			"id":   0,
		},
	})
}
