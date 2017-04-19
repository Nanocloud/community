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

package config

import (
	"github.com/Nanocloud/community/nanocloud/config"
	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
)

func Migrate() error {
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'config'`)
	if err != nil {
		log.Error("Select tables names failed: ", err.Error())
		return err
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("config table already set up")
		return nil
	}
	rows, err = db.Query(
		`CREATE TABLE config (
			key	varchar(255) PRIMARY KEY,
			value varchar(255)
		);`)
	if err != nil {
		log.Errorf("Unable to create config table: %s", err)
		return err
	}

	rows.Close()
	config.Set("windowsAdmin", utils.Env("WINDOWS_USER", "Administrator"), false)
	return nil
}
