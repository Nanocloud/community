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
	UserId       string `json:"user_id"`
	ConnectionId string `json:"connection_id"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
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

	return c.JSON(http.StatusOK, hash{"data": histories})
}

// Add a new log entry to the database
func Add(c *echo.Context) error {
	var t struct {
		Data HistoryInfo
	}
	err := utils.ParseJSONBody(c, &t)
	if err != nil {
		return nil
	}

	if t.Data.UserId == "" || t.Data.ConnectionId == "" || t.Data.StartDate == "" || t.Data.EndDate == "" {
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
		`, t.Data.UserId, t.Data.ConnectionId, t.Data.StartDate, t.Data.EndDate)
	if err != nil {
		return err
	}

	rows.Close()

	return c.JSON(http.StatusCreated, hash{
		"data": hash{
			"success": true,
		},
	})
}

func init() {
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'histories'`)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("Histories table already set up")
		return
	}

	rows, err = db.Query(
		`CREATE TABLE histories (
			userid        varchar(36) NOT NULL DEFAULT '',
			connectionid  varchar(36) NOT NULL DEFAULT '',
			startdate     varchar(36) NOT NULL DEFAULT '',
			enddate       varchar(36) NOT NULL DEFAULT ''
		);`)
	if err != nil {
		log.Errorf("Unable to create histories table: %s", err)
		panic(err)
	}

	rows.Close()
}
