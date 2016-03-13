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

package users

import (
	"errors"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
)

func Migrate() error {
	rows, err := db.Query(
		`SELECT table_name
			FROM information_schema.tables
			WHERE table_name = 'users'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return nil
	}

	rows, err = db.Query(
		`CREATE TABLE users (
				id               varchar(36) PRIMARY KEY,
				firstname        varchar(36) NOT NULL DEFAULT '',
				lastname         varchar(36) NOT NULL DEFAULT '',
				email            varchar(36) NOT NULL DEFAULT '' UNIQUE,
				password         varchar(60) NOT NULL DEFAULT '',
				isadmin          boolean,
				activated        boolean,
				sam              varchar(35) NOT NULL DEFAULT '',
				windowspassword  varchar(36) NOT NULL DEFAULT ''
			);`)
	if err != nil {
		return err
	}

	rows.Close()

	admin, err := users.CreateUser(
		true,
		"admin@nanocloud.com",
		"John",
		"Doe",
		"admin",
		true,
	)

	if err != nil {
		return err
	}

	password := utils.Env("WIN_PASSWORD", "")
	sam := utils.Env("WIN_USER", "")

	result, err := db.Exec(
		`UPDATE users
				SET sam = $1::varchar,
				windows_password = $2::varchar
				WHERE id = $3::varchar;`,
		sam,
		password,
		admin.Id,
	)

	if err != nil {
		log.Error("Failed to update admin account: ", err)
		return err
	}

	updated, err := result.RowsAffected()
	if err != nil {
		log.Error(err)
		return err
	}

	if updated != 1 {
		return errors.New("Unable to set admin password")
	}
	return nil
}
