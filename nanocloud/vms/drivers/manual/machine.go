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
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/Nanocloud/community/nanocloud/vms"
	"github.com/labstack/gommon/log"
)

type machine struct {
	id       string
	server   string
	sshport  string
	user     string
	password string
	role     string
}

func (m *machine) Status() (vms.MachineStatus, error) {
	var service string
	if m.role == "ad" {
		service = "ADWS"
	} else if m.role == "exec" {
		service = "RDMS"
	}
	cmd := exec.Command(
		"sshpass", "-p", m.password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "ConnectTimeout=1",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", m.sshport,
		fmt.Sprintf(
			"%s@%s",
			m.user,
			m.server,
		),
		"powershell.exe \"Write-Host (Get-Service -Name "+service+").status\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		log.Info("Failed to check windows' state", err, string(response))
		return vms.StatusDown, nil
	}
	if strings.Contains(string(response), "Running") {
		return vms.StatusUp, nil
	}
	return vms.StatusDown, nil
}

func (m *machine) IP() (net.IP, error) {
	return []byte(m.server), nil
}

func (m *machine) Type() (vms.MachineType, error) {
	return defaultType, nil
}

func (m *machine) Start() error {
	return nil
}

func (m *machine) Stop() error {
	return nil
}

func (m *machine) Terminate() error {
	return nil
}

func (m *machine) Id() string {
	return m.id
}

func (m *machine) Name() (string, error) {
	if m.role == "ad" {
		return "Windows Active Directory", nil
	} else if m.role == "exec" {
		return "Windows Session Host", nil
	}
	return "Undefined Windows Server", nil
}
