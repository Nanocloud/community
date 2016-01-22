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
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/Nanocloud/nano"
	_ "github.com/lib/pq"
)

var module nano.Module
var db *sql.DB

type hash map[string]interface{}

// Log entries are stored in this structure
type HistoryInfo struct {
	UserId       string
	ConnectionId string
	StartDate    string
	EndDate      string
}

func dbConnect() {
	databaseURI := os.Getenv("DATABASE_URI")
	if len(databaseURI) == 0 {
		databaseURI = "postgres://localhost/nanocloud?sslmode=disable"
	}

	var err error

	for try := 0; try < 10; try++ {
		db, err = sql.Open("postgres", databaseURI)
		if err == nil {
			return
		}
		time.Sleep(time.Second * 5)
	}

	module.Log.Fatalf("Cannot connect to Postgres Database: %s", err)
}

// Connects to the postgres databse
func setupDb() error {
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'histories'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		module.Log.Info("Histories table already set up")
		return nil
	}

	rows, err = db.Query(
		`CREATE TABLE histories (
			userid		varchar(36),
			connectionid	varchar(36),
			startdate	varchar(36),
			enddate		varchar(36)
		);`)
	if err != nil {
		module.Log.Errorf("Unable to create histories table: %s", err)
		return err
	}

	rows.Close()
	return nil
}

// Get a list of all the log entries of the database
func ListCall(req nano.Request) (*nano.Response, error) {
	var histories []HistoryInfo
	rows, err := db.Query(
		`SELECT userid, connectionid,
		startdate, enddate
		FROM histories`,
	)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	if len(histories) == 0 {
		histories = []HistoryInfo{}
	}

	return nano.JSONResponse(200, histories), nil
}

// Add a new log entry to the database
func AddCall(req nano.Request) (*nano.Response, error) {
	var t HistoryInfo
	err := json.Unmarshal([]byte(req.Body), &t)
	if err != nil {
		return nano.JSONResponse(400, hash{
			"error": "Invalid parameters",
		}), nil
	}

	rows, err := db.Query(
		`INSERT INTO histories
		(userid, connectionid, startdate, enddate)
		VALUES(	$1::varchar, $2::varchar, $3::varchar, $4::varchar)
		`, t.UserId, t.ConnectionId, t.StartDate, t.EndDate)
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	rows.Close()

	return nano.JSONResponse(201, hash{
		"success": true,
	}), nil
}

func main() {
	module = nano.RegisterModule("history")

	dbConnect()
	setupDb()

	module.Get("/history", ListCall)
	module.Post("/history", AddCall)

	module.Listen()
}
