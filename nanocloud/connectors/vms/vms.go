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

package vms

import (
	"os"

	"github.com/Nanocloud/community/nanocloud/vms"
	_ "github.com/Nanocloud/community/nanocloud/vms/drivers/qemu"
	log "github.com/Sirupsen/logrus"
)

var _vm *vms.VM

func getInstance() (*vms.VM, error) {
	if _vm == nil {

		iaas := os.Getenv("IAAS")
		if len(iaas) == 0 {
			log.Fatal("No iaas provided")
		}
		server := os.Getenv("WIN_SERVER")
		m := make(map[string]string, 1)
		m["server"] = server
		var err error
		_vm, err = vms.Open(iaas, m)
		return _vm, err
	}
	return _vm, nil
}

func Machines() ([]vms.Machine, error) {
	vm, err := getInstance()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return (*vm).Machines()
}

func Machine(id string) (vms.Machine, error) {
	vm, err := getInstance()
	if err != nil {
		return nil, err
	}
	return (*vm).Machine(id)
}

func Create(name, password string, t vms.MachineType) (vms.Machine, error) {
	vm, err := getInstance()
	if err != nil {
		return nil, err
	}
	return (*vm).Create(name, password, t)
}
