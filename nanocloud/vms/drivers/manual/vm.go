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

package manual

import (
	"fmt"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/vms"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
)

type vm struct {
}

func (v *vm) Types() ([]vms.MachineType, error) {
	return []vms.MachineType{defaultType}, nil
}

func (v *vm) Create(attr vms.MachineAttributes) (vms.Machine, error) {

	machine := &machine{
		id:       uuid.NewV4().String(),
		server:   attr.Ip,
		user:     attr.Username,
		password: attr.Password,
	}
	rows, err := db.Query(
		`INSERT INTO machines
		(id, type, ad, execserver, username, password)
		VALUES( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::varchar, $6::varchar)`,
		machine.id, "manual", machine.server, machine.server, machine.user, machine.password)
	if err != nil {
		return nil, err
	}
	rows.Close()
	return machine, nil
}

func (v *vm) Machines() ([]vms.Machine, error) {

	rows, err := db.Query(
		`SELECT role, type, id, execserver,
		plazaport, username, password
		FROM machines`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	machines := make([]vms.Machine, 0)
	vmType := ""
	for rows.Next() {
		machine := &machine{}

		rows.Scan(
			&machine.role,
			&vmType,
			&machine.id,
			&machine.server,
			&machine.plazaport,
			&machine.user,
			&machine.password,
		)
		if vmType == "manual" {
			machines = append(machines, machine)
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return machines, nil
}

func (v *vm) Machine(id string) (vms.Machine, error) {
	rows, err := db.Query(
		`SELECT role, id, execserver, plazaport, username, password
		FROM machines WHERE id = $1::varchar`, id)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	machine := &machine{}
	for rows.Next() {
		rows.Scan(
			&machine.role,
			&machine.id,
			&machine.server,
			&machine.plazaport,
			&machine.user,
			&machine.password,
		)
	}
	err = rows.Err()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return machine, nil
}
