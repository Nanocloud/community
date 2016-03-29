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

import "net"

type MachineStatus int

const (
	StatusUnknown    MachineStatus = 0
	StatusDown       MachineStatus = 1
	StatusUp         MachineStatus = 2
	StatusTerminated MachineStatus = 3
	StatusBooting    MachineStatus = 4
	StatusCreating   MachineStatus = 5
)

type Machine interface {
	Id() string
	Platform() string
	Name() (string, error)
	Status() (MachineStatus, error)
	IP() (net.IP, error)
	Type() (MachineType, error)
	Progress() (uint8, error)

	Start() error
	Stop() error
	Terminate() error
}

func StatusToString(status MachineStatus) string {
	switch status {
	case StatusDown:
		return "down"
	case StatusUp:
		return "up"
	case StatusTerminated:
		return "terminated"
	case StatusBooting:
		return "booting"
	case StatusCreating:
		return "creating"
	}
	return "unknown"
}
