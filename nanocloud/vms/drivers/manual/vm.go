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
	"errors"
	"strings"

	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type vm struct {
	servers  string
	ad       string
	sshport  string
	user     string
	password string
}

func (v *vm) Types() ([]vms.MachineType, error) {
	return []vms.MachineType{defaultType}, nil
}

func (v *vm) Create(name, password string, t vms.MachineType) (vms.Machine, error) {
	log.Error("Not Implemented")
	return nil, nil
}

func (v *vm) Machines() ([]vms.Machine, error) {
	if v.ad == v.servers {
		machines := make([]vms.Machine, 1)
		machines[0] = &machine{role: "ad", id: v.ad, server: v.ad, sshport: v.sshport, user: v.user, password: v.password}
		return machines, nil
	}
	ips := strings.Split(v.servers, ";")
	var machines []vms.Machine
	machines = append(machines, &machine{role: "ad", id: v.ad, server: v.ad, sshport: v.sshport, user: v.user, password: v.password})
	for _, val := range ips {
		machines = append(machines, &machine{role: "exec", id: val, server: val, sshport: v.sshport, user: v.user, password: v.password})
	}
	return machines, nil
}

func (v *vm) Machine(id string) (vms.Machine, error) {
	return nil, errors.New("Not implemented")
}
