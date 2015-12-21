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
	"errors"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"regexp"

	"github.com/natefinch/pie"
)

var (
	name = "history"
	srv  pie.Server
)

type HistoryConfig struct {
	ConnectionString string
	DatabaseName     string
}

// Plugin Structure
type api struct{}

var db *sql.DB

// Log entries are stored in this structure
type HistoryInfo struct {
	UserId       string
	ConnectionId string
	StartDate    string
	EndDate      string
}

// Structure used for exchanges with the core, faking a request/responsewriter
type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
	Method   string
	HeadVals map[string]string
	Status   int
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
		log.Info("[History] histories table already set up\n")
		return nil
	}

	rows, err = db.Query(
		`CREATE TABLE histories (
			userid		varchar(36),
			connectionid	varchar(36) PRIMARY KEY,
			startdate	varchar(36),
			enddate		varchar(36)
		);`)
	if err != nil {
		log.Errorf("[History] Unable to create histories table: %s\n", err)
		return err
	}

	rows.Close()
	return nil
}

// Get a list of all the log entries of the database
func GetList(histories *[]HistoryInfo) error {
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
		*histories = append(*histories, history)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil

}

// Add a new log entry to the database
func Add(args PlugRequest, reply *PlugRequest) error {

	var t HistoryInfo
	err := json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		reply.Status = 400
		log.Error("Json Arguments are not valid: ", err)
		return err
	}
	rows, err := db.Query(
		`INSERT INTO histories
		(userid, connectionid, startdate, enddate)
		VALUES(	$1::varchar, $2::varchar, $3::varchar, $4::varchar)
		`, t.UserId, t.ConnectionId, t.StartDate, t.EndDate)
	if err != nil {
		switch err.Error() {
		case "pq: duplicate key value violates unique constraint \"history_pkey\"":
			err = errors.New("history connection id exists already")
		}
		return err
	}

	rows.Close()
	return err

}

// Connects to the DB, adds entry ,sets the status and body of the response and closes the DB
func AddCall(args PlugRequest, reply *PlugRequest, id string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	initConf()
	err := Add(args, reply)
	if err != nil {
		if reply.Status != 400 {
			reply.Status = 500
		}
		log.Error("Failed to add an entry: ", err)
		return
	}
	reply.Status = 202
}

// Connects to the DB, list entries ,sets the status and body of the response and closes the DB
func ListCall(args PlugRequest, reply *PlugRequest, id string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	initConf()
	var histories []HistoryInfo
	GetList(&histories)
	rsp, err := json.Marshal(histories)
	if err != nil {
		reply.Status = 500
		log.Error("Failed to Marshal histories: ", err)
		return
	}
	reply.Body = string(rsp)
	reply.Status = 200

}

// slice of available handler functions
var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string)
}{
	{`^\/api\/history\/{0,1}$`, "GET", ListCall},
	{`^\/api\/history\/{0,1}$`, "POST", AddCall},
}

// Will receive all http requests starting by /api/history from the core and chose the correct handler function
func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
				return nil
			} else {
				val.f(args, reply, "")
				return nil
			}
		}
	}
	return nil
}

// Plug the plugin to the core
func (api) Plug(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

// Will contain various verifications for the plugin. If the core can call the function and receives "true" in the reply, it means the plugin is functionning correctly
func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

// Unplug the plugin from the core
func (api) Unplug(args interface{}, reply *bool) error {
	defer os.Exit(0)
	*reply = true
	return nil
}

func main() {

	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	initConf()
	srv = pie.NewProvider()
	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatal("Failed to register:", name, err)
	}

	var err error
	db, err = sql.Open("postgres", conf.DatabaseUri)
	if err != nil {
		log.Fatalf("Cannot connect to Postgres Database: %s", err)
	}

	err = setupDb()
	if err != nil {
		log.Fatalf("[History] unable to setup histories table: %s", err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)
}
