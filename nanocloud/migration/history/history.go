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
	"github.com/Nanocloud/community/nanocloud/connectors/db"
	log "github.com/Sirupsen/logrus"
)

func Migrate() error {
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
		return nil
	}

	rows, err = db.Query(
		`CREATE TABLE histories (
			id			varchar(36) PRIMARY KEY,
			userid			varchar(36) NOT NULL DEFAULT '',
			usermail		varchar(36) NOT NULL DEFAULT '',
			userfirstname	varchar(36) NOT NULL DEFAULT '',
			userlastname	varchar(36) NOT NULL DEFAULT '',
			connectionid	varchar(36) NOT NULL DEFAULT '',
			startdate		varchar(36) NOT NULL DEFAULT '',
			enddate			varchar(36) NOT NULL DEFAULT ''
		);`)
	if err != nil {
		log.Errorf("Unable to create histories table: %s", err)
		return err
	}

	rows.Close()
	return nil
}
