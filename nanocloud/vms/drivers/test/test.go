/* Nanocloud Community, a comprehensive platform to turn any application
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

package test

import (
	"errors"
	"net"
	"time"

	"github.com/Nanocloud/community/nanocloud/vms"
	uuid "github.com/satori/go.uuid"
)

type machineType struct {
	flavour string
}

func (t *machineType) GetID() string {
	return "default-test-machine-type"
}

var (
	allMachines        []vms.Machine
	defaultMachineType machineType
	failNextCall       bool  = false
	nilNextCall        bool  = false
	waitNextCall       bool  = false
	waitDelay          int   = 0
	failError          error = errors.New("fake-generated-error")
)

type machine struct {
	id          string
	name        string
	status      vms.MachineStatus
	machineType machineType
}

func SetNil() {
	nilNextCall = true
}

func SetFail() {
	failNextCall = true
}

func SetDelay(delay int) {
	waitNextCall = true
	waitDelay = delay
}

func (m *machine) Id() string {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if nilNextCall {
		nilNextCall = false
		return ""
	}
	return m.id
}

func (m *machine) Platform() string {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if nilNextCall {
		nilNextCall = false
		return ""
	}
	return "test"
}

func (m *machine) Name() (string, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return "", failError
	}
	if nilNextCall {
		nilNextCall = false
		return "", nil
	}
	return m.name, nil
}

func (m *machine) Status() (vms.MachineStatus, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return m.status, failError
	}
	if nilNextCall {
		nilNextCall = false
		return m.status, nil
	}
	return m.status, nil
}

func (m *machine) IP() (net.IP, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return nil, failError
	}
	if nilNextCall {
		nilNextCall = false
		return nil, nil
	}
	return net.ParseIP("127.0.0.1"), nil
}

func (m *machine) Type() (vms.MachineType, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return nil, failError
	}
	if nilNextCall {
		nilNextCall = false
		return nil, nil
	}
	return &defaultMachineType, nil
}

func (m *machine) Progress() (uint8, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return 0, failError
	}
	return 0, nil
}

func (m *machine) Credentials() (string, string, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return "", "", failError
	}
	if nilNextCall {
		nilNextCall = false
		return "", "", nil
	}
	return "credential-test", "credential-test", nil
}

func (m *machine) Start() error {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return failError
	}
	m.status = vms.StatusUp
	return nil
}

func (m *machine) Stop() error {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return failError
	}
	m.status = vms.StatusStopping

	go func() {
		time.Sleep(1000 * time.Millisecond)
		m.status = vms.StatusDown
	}()
	return nil
}

func (m *machine) Terminate() error {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return failError
	}
	m.status = vms.StatusTerminated
	return nil
}

type vm struct {
}

func (v *vm) Machines() ([]vms.Machine, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return nil, failError
	}
	return allMachines, nil
}

func (v *vm) Machine(id string) (vms.Machine, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return nil, failError
	}
	if nilNextCall {
		nilNextCall = false
		return nil, nil
	}
	for _, machine := range allMachines {
		if machine.Id() == id {
			return machine, nil
		}
	}
	return nil, nil
}

func (v *vm) Create(attr vms.MachineAttributes) (vms.Machine, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return nil, failError
	}
	if nilNextCall {
		nilNextCall = false
		return nil, nil
	}
	if attr.Type == nil {
		attr.Type = &defaultMachineType
	}

	t, ok := attr.Type.(*machineType)
	if !ok {
		return nil, errors.New("VM Type not supported")
	}

	new_machine := machine{
		id:          uuid.NewV4().String(),
		name:        attr.Name,
		machineType: *t,
		status:      vms.StatusUp,
	}
	allMachines = append(allMachines, &new_machine)

	return &new_machine, nil
}

func (v *vm) Types() ([]vms.MachineType, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return nil, failError
	}
	if nilNextCall {
		nilNextCall = false
		return nil, nil
	}
	rt := make([]vms.MachineType, 1)
	rt[0] = &defaultMachineType
	return rt, nil
}

func (v *vm) Type(id string) (vms.MachineType, error) {
	if waitNextCall {
		time.Sleep(time.Duration(waitDelay) * time.Millisecond)
		waitNextCall = false
		waitDelay = 0
	}
	if failNextCall {
		failNextCall = false
		return nil, failError
	}
	if nilNextCall {
		nilNextCall = false
		return nil, nil
	}
	if id == defaultMachineType.GetID() {
		return &defaultMachineType, nil
	}
	return nil, errors.New("Type not found")
}

type driver struct {
}

func (d *driver) Open(options map[string]string) (vms.VM, error) {
	defaultMachineType.flavour = "t2.tiny"
	return &vm{}, nil
}

func init() {
	vms.Register("test", &driver{})
}
