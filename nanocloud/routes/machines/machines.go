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

package machines

import (
	"net/http"

	"github.com/Nanocloud/community/nanocloud/connectors/vms"
	"github.com/Nanocloud/community/nanocloud/utils"
	vm "github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

type sMachine struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name          string `json:"name"`
		Ip            string `json:"ip"`
		Status        string `json:"status"`
		AdminPassword string `json:"admin-password,omitempty"`
	} `json:"attributes"`
}

func serializableMachine(m vm.Machine) *sMachine {
	rt := sMachine{}
	rt.Id = m.Id()
	rt.Type = "machine"
	rt.Attributes.Name, _ = m.Name()
	status, _ := m.Status()
	rt.Attributes.Status = vm.StatusToString(status)
	ip, _ := m.IP()
	rt.Attributes.Ip = string(ip)
	return &rt
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

func PatchMachine(c *echo.Context) error {
	body := struct {
		Data sMachine `json:"data"`
	}{}

	err := utils.ParseJSONBody(c, &body)
	if err != nil {
		return err
	}

	m, err := vms.Machine(body.Data.Id)
	if err != nil {
		return err
	}

	s, err := m.Status()
	if err != nil {
		return err
	}

	if body.Data.Attributes.Status == "up" {
		if s == vms.StatusDown {
			err = m.Start()
			if err != nil {
				return err
			}
		}
	} else if body.Data.Attributes.Status == "down" {
		if s == vms.StatusUp {
			err = m.Stop()
			if err != nil {
				return err
			}
		}
	}

	m, err = vms.Machine(m.Id())
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, hash{"data": serializableMachine(m)})
}

func GetMachine(c *echo.Context) error {
	m, err := vms.Machine(c.Param("id"))
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

	var machine struct {
		Id   string `json:"id"`
		Type string `json:"type"`
		Attr struct {
			Name   string `json:"name"`
			Ip     string `json:"ip"`
			Status string `json:"status"`
		} `json:"attributes"`
	}

	machine.Id = m.Id()
	machine.Type = "machine"
	machine.Attr.Name, err = m.Name()
	if err != nil {
		log.Error(err)
		return retJsonError(c, err)
	}
	status, err := m.Status()
	if err != nil {
		log.Error(err)
		return retJsonError(c, err)
	}
	machine.Attr.Status = vm.StatusToString(status)

	ip, err := m.IP()
	machine.Attr.Ip = ip.String()
	if err != nil {
		log.Error(err)
		return retJsonError(c, err)
	}

	return c.JSON(http.StatusOK, hash{"data": machine})
}

func Machines(c *echo.Context) error {
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
		Name   string `json:"name"`
		Ip     string `json:"ip"`
		Status string `json:"status"`
	}
	type virtmachine struct {
		Id   string `json:"id"`
		Type string `json:"type"`
		Att  attr   `json:"attributes"`
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
		res[i].Type = "machine"
		res[i].Att.Status = vm.StatusToString(status)
		ip, _ := val.IP()
		if ip != nil {
			res[i].Att.Ip = ip.String()
		}
	}

	return c.JSON(http.StatusOK, hash{"data": res})
}

func CreateMachine(c *echo.Context) error {
	body := struct {
		Data sMachine `json:"data"`
	}{}

	err := utils.ParseJSONBody(c, &body)
	if err != nil {
		log.Error(err)
		return retJsonError(c, err)
	}

	vm, err := vms.Create(body.Data.Attributes.Name, body.Data.Attributes.AdminPassword, nil)
	if err != nil {
		log.Error(err)
		return retJsonError(c, err)
	}
	return c.JSON(http.StatusOK, hash{
		"data": serializableMachine(vm),
	})
}

func DeleteMachine(c *echo.Context) error {
	id := c.Param("id")

	machine, err := vms.Machine(id)
	if err != nil {
		return err
	}
	err = machine.Terminate()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, hash{})
}
