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
	"github.com/Nanocloud/community/nanocloud/errors"
	"github.com/Nanocloud/community/nanocloud/utils"
	vm "github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

type hash map[string]interface{}

type machine struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Ip            string `json:"ip"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	UserName      string `json:"username"`
	AdminPassword string `json:"admin-password,omitempty"`
	Platform      string `json:"platform"`
	Progress      int    `json:"progress"`
}

func (m *machine) GetID() string {
	return m.Id
}

func (m *machine) SetID(id string) error {
	m.Id = id
	return nil
}

func PatchMachine(c *echo.Context) error {
	b := &machine{}

	err := utils.ParseJSONBody(c, b)
	if err != nil {
		log.Error(err)
		return errors.UnableToUpdateMachineStatus
	}

	m, err := vms.Machine(b.Id)
	if err != nil {
		log.Error(err)
		return errors.UnableToUpdateMachineStatus
	}

	status, err := m.Status()
	if err != nil {
		log.Error(err)
		return errors.UnableToUpdateMachineStatus
	}

	switch b.Status {
	case "up":
		if status != vms.StatusDown {
			return errors.UnableToUpdateMachineStatus
		}

		err = m.Start()
		if err != nil {
			log.Error(err)
			return errors.UnableToUpdateMachineStatus
		}

	case "down":
		if status != vms.StatusUp {
			return errors.UnableToUpdateMachineStatus
		}
		err = m.Stop()
		if err != nil {
			log.Error(err)
			return errors.UnableToUpdateMachineStatus
		}

	default:
		log.Error(err)
		return errors.UnableToUpdateMachineStatus
	}

	m, err = vms.Machine(m.Id())
	if err != nil {
		log.Error(err)
		return errors.UnableToUpdateMachineStatus
	}

	rt, err := getSerializableMachine(m.Id())
	if err != nil {
		log.Error(err)
		return errors.UnableToUpdateMachineStatus
	}

	return utils.JSON(c, http.StatusOK, rt)
}

func getSerializableMachine(id string) (*machine, error) {
	m, err := vms.Machine(id)
	if err != nil {
		return nil, err
	}

	rt := machine{}

	rt.Id = m.Id()
	rt.Platform = m.Platform()

	rt.Name, err = m.Name()
	if err != nil {
		return nil, err
	}

	status, err := m.Status()
	if err != nil {
		return nil, err
	}
	rt.Status = vm.StatusToString(status)

	ip, err := m.IP()
	if err != nil {
		return nil, err
	}

	if ip != nil {
		rt.Ip = ip.String()
	}

	return &rt, nil
}

func GetMachine(c *echo.Context) error {
	m, err := getSerializableMachine(c.Param("id"))
	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusOK, m)
}

func Machines(c *echo.Context) error {
	machines, err := vms.Machines()
	if err != nil {
		log.Error(err)
		return errors.UnableToRetrieveMachineList
	}

	res := make([]*machine, len(machines))

	for i, val := range machines {
		m := machine{}
		m.Name, err = val.Name()
		if err != nil {
			log.Error(err)
			return errors.UnableToRetrieveMachineList
		}

		m.Id = val.Id()
		status, err := val.Status()
		if err != nil {
			log.Error(err)
			return errors.UnableToRetrieveMachineList
		}

		m.Platform = val.Platform()

		m.Status = vm.StatusToString(status)

		progress, err := val.Progress()
		if err != nil {
			log.Errorf("Unable to get machine progress: %s", err)
		} else {
			m.Progress = int(progress)
		}

		ip, _ := val.IP()
		if ip != nil {
			m.Ip = ip.String()
		}

		res[i] = &m
	}

	return utils.JSON(c, http.StatusOK, res)
}

func CreateMachine(c *echo.Context) error {
	rt := &machine{}

	err := utils.ParseJSONBody(c, rt)
	if err != nil {
		return err
	}

	attr := vm.MachineAttributes{
		Type:     nil,
		Name:     rt.Name,
		Username: rt.UserName,
		Password: rt.AdminPassword,
		Ip:       rt.Ip,
	}
	m, err := vms.Create(attr)
	if err != nil {
		log.Error(err)
		return errors.UnableToCreateTheMachine
	}

	rt, err = getSerializableMachine(m.Id())
	if err != nil {
		return err
	}

	return utils.JSON(c, http.StatusOK, rt)
}

func DeleteMachine(c *echo.Context) error {
	id := c.Param("id")

	m, err := vms.Machine(id)
	if err != nil {
		log.Error(err)
		return errors.UnableToTerminateTheMachine
	}

	err = m.Terminate()
	if err != nil {
		log.Error(err)
		return errors.UnableToTerminateTheMachine
	}
	return c.JSON(http.StatusOK, hash{"meta": hash{}})
}
