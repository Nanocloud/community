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

package machines

import (
	"fmt"
	"os"
	"strings"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

func Migrate() error {
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'machines'`)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("Machines table already set up")
		return nil
	}

	_, err = db.Exec(`CREATE TYPE vmtype AS ENUM('qemu', 'manual', 'aws', 'vmware', 'azure')`)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	rows, err = db.Query(
		`CREATE TABLE machines (
			id         varchar(60),
			name			 varchar(255),
			type       vmtype,
			ip         varchar(255),
			plazaport  varchar(4) NOT NULL DEFAULT '9090',
			username   varchar(36),
			password   varchar(60)
		);`)
	if err != nil {
		log.Errorf("Unable to create machines table: %s", err)
		return err
	}
	if os.Getenv("IAAS") == "manual" {
		servers := os.Getenv("EXECUTION_SERVERS")
		password := os.Getenv("WINDOWS_PASSWORD")
		user := os.Getenv("WINDOWS_USER")

		ips := strings.Split(servers, ";")
		for i, val := range ips {
			rows, err := db.Query(`INSERT INTO machines
			(id, name, type, ip, username, password)
			VALUES( $1::varchar, $2::varchar, $3::vmtype,
			$4::varchar, $5::varchar, $6::varchar)`,
				uuid.NewV4().String(),
				fmt.Sprintf("Machine #%d", i+1),
				"manual",
				val,
				user,
				password,
			)

			if err != nil {
				return err
			}
			rows.Close()
		}
	}
	rows.Close()
	return nil
}
