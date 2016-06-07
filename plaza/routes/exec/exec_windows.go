// +build windows

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

package exec

import (
	"bytes"

	"github.com/Nanocloud/community/plaza/windows"
)

func runCommand(username string, domain string, password string, command []string) windows.Cmd {
	cmd := windows.Command(
		username, domain, password,
		command[0], command[1:]...,
	)
	return *cmd
}

func launchApp(command []string) (uint32, error) {
	for tries := 20; tries > 0; tries-- {
		pid, err = windows.LaunchApp(body.Command)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}

	if err != nil {
		return 0, err
	}

	return pid, err
}

func makeResponse(stdout bytes.Buffer, stderr bytes.Buffer, cmd windows.Cmd) map[string]interface{} {
	res := make(map[string]interface{})
	res["stdout"] = stdout.String()
	res["stderr"] = stderr.String()
	res["time"] = cmd.ProcessState.SysUsage()
	res["code"] = cmd.ProcessState.String()

	return res
}
