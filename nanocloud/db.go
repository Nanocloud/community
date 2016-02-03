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
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"time"
)

var db *sql.DB = nil

func dbConnect() error {
	var err error
	uri := env("DATABASE_URI", "postgres://localhost/nanocloud?sslmode=disable")

	for try := 0; try < 10; try++ {
		db, err = sql.Open("postgres", uri)
		if err != nil {
			return err
		}

		err = db.Ping()
		if err == nil {
			module.Log.Info("Connected to Postgres")
			return nil
		}

		module.Log.Info("Unable to connect to Postgres. Will retry in 5 sec")
		time.Sleep(time.Second * 5)
	}
	return err
}

func setupDb() error {
	// oauth_clients table
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'oauth_clients'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("[nanocloud] oauth_clients table already set up\n")
	} else {
		rows, err = db.Query(
			`CREATE TABLE oauth_clients (
				id      serial PRIMARY KEY,
				name    varchar(255) NOT NULL DEFAULT '' UNIQUE,
				key     varchar(255) NOT NULL DEFAULT '' UNIQUE,
				secret  varchar(255) NOT NULL DEFAULT ''
			)`)

		if err != nil {
			log.Errorf("[nanocloud] Unable to create oauth_clients table: %s\n", err)
			return err
		}
		defer rows.Close()

		rows, err = db.Query(
			`INSERT INTO oauth_clients
			(name, key, secret)
			VALUES (
				'Nanocloud',
				'9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae',
				'9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341'
			)`)

		if err != nil {
			log.Errorf("[nanocloud] Unable to create default oauth_clients: %s\n", err)
			return err
		}
		defer rows.Close()
	}

	// oauth_access_tokens table
	rows, err = db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'oauth_access_tokens'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("[nanocloud] oauth_access_tokens table already set up\n")
	} else {
		rows, err = db.Query(
			`CREATE TABLE oauth_access_tokens (
				id                serial PRIMARY KEY,
				token             varchar(255) NOT NULL DEFAULT '' UNIQUE,
				oauth_client_id   integer REFERENCES oauth_clients (id),
				user_id           varchar(255) NOT NULL DEFAULT ''
			)`)

		if err != nil {
			log.Errorf("[nanocloud] Unable to create oauth_access_tokens table: %s\n", err)
			return err
		}
		defer rows.Close()
	}
	return nil
}
