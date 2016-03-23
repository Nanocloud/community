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

type MachineStatus struct {
	Status      int    `json:"status"`
	CurrentSize string `json:"current-size"`
	TotalSize   string `json:"total-size"`
}

var (
	StatusUnknown     MachineStatus = MachineStatus{Status: 0}
	StatusDown        MachineStatus = MachineStatus{Status: 1}
	StatusUp          MachineStatus = MachineStatus{Status: 2}
	StatusTerminated  MachineStatus = MachineStatus{Status: 3}
	StatusBooting     MachineStatus = MachineStatus{Status: 4}
	StatusDownloading MachineStatus = MachineStatus{Status: 5}
)

type Machine interface {
	Id() string
	Platform() string
	Name() (string, error)
	Status() (MachineStatus, error)
	IP() (net.IP, error)
	Type() (MachineType, error)

	Start() error
	Stop() error
	Terminate() error
}

func StatusToString(status MachineStatus) string {
	switch status.Status {
	case 1:
		return "available"
	case 2:
		return "running"
	case 3:
		return "terminated"
	case 4:
		return "booting"
	case 5:
		return "download"
	}
	return "unknown"
}
