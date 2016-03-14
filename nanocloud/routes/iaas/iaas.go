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

package iaas

import (
	"net/http"

	"github.com/Nanocloud/community/nanocloud/connectors/vms"
	vm "github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

type jsonMachine struct {
	Id     string
	Name   string
	Status string
	Ip     string
}

func machinetoStruct(rawmachine vm.Machine) jsonMachine {
	var mach jsonMachine
	mach.Id = rawmachine.Id()
	mach.Name, _ = rawmachine.Name()
	status, _ := rawmachine.Status()
	mach.Status = vm.StatusToString(status)
	ip, _ := rawmachine.IP()
	mach.Ip = string(ip)
	return mach
}

func retJsonError(c *echo.Context, err error) error {
	return c.JSON(
		http.StatusInternalServerError, hash{
			"errors": [1]hash{
				hash{
					"detail": err.Error(),
				},
			},
		})
}

func ListRunningVM(c *echo.Context) error {
	machines, err := vms.Machines()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError, hash{
				"errors": [1]hash{
					hash{
						"detail": err.Error(),
					},
				},
			})
	}
	type attr struct {
		Name        string `json:"name"`
		Ip          string `json:"ip"`
		Status      string `json:"status"`
		TotalSize   string `json:"total-size"`
		CurrentSize string `json:"current-size"`
	}
	type virtmachine struct {
		Id  string `json:"id"`
		Att attr   `json:"attributes"`
	}
	var res = make([]virtmachine, len(machines))
	for i, val := range machines {
		res[i].Att.Name, err = val.Name()
		if err != nil {
			log.Error(err)
			return retJsonError(c, err)
		}
		res[i].Id = val.Id()
		status, err := val.Status()
		if err != nil {
			log.Error(err)
			return retJsonError(c, err)
		}
		res[i].Att.Status = vm.StatusToString(status)
		res[i].Att.CurrentSize = status.CurrentSize
		res[i].Att.TotalSize = status.TotalSize
		ip, _ := val.IP()
		res[i].Att.Ip = ip.String()
		if err != nil {
			log.Error(err)
			return retJsonError(c, err)
		}
	}

	return c.JSON(http.StatusOK, hash{"data": res})
}

func StopVM(c *echo.Context) error {
	machine, err := vms.Machine(c.Param("id"))

	err = machine.Stop()
	if err != nil {
		return retJsonError(c, err)
	}
	return c.JSON(
		http.StatusOK, hash{
			"vm": machinetoStruct(machine),
		})
}

func StartVM(c *echo.Context) error {
	machine, err := vms.Machine(c.Param("id"))
	if err != nil {
		return retJsonError(c, err)
	}

	err = machine.Start()
	if err != nil {
		return retJsonError(c, err)
	}
	return c.JSON(
		http.StatusOK, hash{
			"vm": machinetoStruct(machine),
		})
}

func CreateVM(c *echo.Context) error {
	//TODO READ BODY TO GET PASSWORD AND TYPE
	vm, err := vms.Create(c.Param("id"), "", nil)
	if err != nil {
		return retJsonError(c, err)
	}
	return c.JSON(
		http.StatusOK, hash{
			"vm": machinetoStruct(vm),
		})
}
