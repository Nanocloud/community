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

package qemu

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type vm struct {
	server string
}

func (v *vm) Types() ([]vms.MachineType, error) {
	return []vms.MachineType{defaultType}, nil
}

func (v *vm) Create(attr vms.MachineAttributes) (vms.Machine, error) {

	m := machine{id: attr.Username, server: v.server}
	ip, _ := m.IP()
	resp, err := http.Post("http://"+string(ip)+":8080/api/vms/"+m.Id()+"/download", "", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Error(err)
		return nil, err
	}
	return &m, nil
}

func (v *vm) Machines() ([]vms.Machine, error) {
	resp, err := http.Get("http://" + v.server + ":8080/api/vms")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	type Status struct {
		Id         string `json:"id"`
		Type       string `json:"type"`
		Attributes vmInfo
	}
	var State struct {
		Data []Status `json:"data"`
	}

	err = json.Unmarshal(body, &State)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var machines = make([]vms.Machine, len(State.Data))
	for i, val := range State.Data {
		machines[i] = &machine{id: val.Id, server: v.server}
	}
	return machines, nil
}

func (v *vm) Machine(id string) (vms.Machine, error) {
	return &machine{id: id, server: v.server}, nil
}
