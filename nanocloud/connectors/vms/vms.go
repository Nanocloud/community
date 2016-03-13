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
	"github.com/Nanocloud/community/nanocloud/vms"
)

const (
	StatusUnknown     = vms.StatusUnknown
	StatusDown        = vms.StatusDown
	StatusUp          = vms.StatusUp
	StatusTerminated  = vms.StatusTerminated
	StatusBooting     = vms.StatusBooting
	StatusDownloading = vms.StatusDownloading
)

var vm *vms.VM

func SetVM(v *vms.VM) {
	vm = v
}

func Machines() ([]vms.Machine, error) {
	return (*vm).Machines()
}

func Machine(id string) (vms.Machine, error) {
	return (*vm).Machine(id)
}

func Create(name, password string, t vms.MachineType) (vms.Machine, error) {
	return (*vm).Create(name, password, t)
}
