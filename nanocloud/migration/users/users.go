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
	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
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
				first_name       varchar(36) NOT NULL DEFAULT '',
				last_name        varchar(36) NOT NULL DEFAULT '',
				email            varchar(36) NOT NULL DEFAULT '' UNIQUE,
				password         varchar(60) NOT NULL DEFAULT '',
				is_admin         boolean,
				activated        boolean,
				sam              varchar(35) NOT NULL DEFAULT '',
				windows_password varchar(36) NOT NULL DEFAULT ''
			);`)
	if err != nil {
		return err
	}

	rows.Close()

	_, err = users.CreateUser(
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
	return nil
}
