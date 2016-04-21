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

package plaza

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type cmd_t struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Domain   string   `json:"domain"`
	Command  []string `json:"command"`
	Stdin    string   `json:"stdin"`
}

type result_t struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Code   uint64 `json:"code"`
	// Time   uint64 `json:"time"`
}

func Exec(address string, port int, cmd *cmd_t) (*result_t, error) {
	var err error

	client := &http.Client{}

	instr, err := json.Marshal(&cmd)
	if err != nil {
		return nil, err
	}

	buff := bytes.NewBuffer(instr)

	var resp *http.Response

	for i := 0; i < 10; i++ {
		resp, err = client.Post(
			fmt.Sprintf("http://%s:%d/exec", address, port),
			"application/json",
			buff,
		)
		if err == nil {
			break
		} else {
			time.Sleep(5 * time.Second)
		}
	}
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	r := result_t{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func PowershellExec(
	address string, port int,
	username string, domain string, password string,
	command ...string,
) (*result_t, error) {
	cmd := cmd_t{
		Username: username,
		Password: password,
		Domain:   domain,
		Command: []string{
			"C:\\Windows\\System32\\WindowsPowershell\\v1.0\\powershell.exe",
			"-Command",
			"-",
		},
		Stdin: strings.Join(command, " "),
	}

	res, err := Exec(address, port, &cmd)
	if err != nil {
		return nil, err
	}

	if res.Code != 0 {
		return nil, errors.New(res.Stdout)
	}
	return res, nil
}

func PublishApp(
	address string, port int,
	username string, domain string, password string,
	collectionName string, displayName string, filePath string,
) ([]byte, error) {
	res, err := PowershellExec(
		address, port,
		username, domain, password,
		"Try {",
		"Import-module RemoteDesktop;",
		fmt.Sprintf(
			"New-RDRemoteApp -CollectionName '%s' -DisplayName '%s' -FilePath '%s' -ErrorAction Stop | ConvertTo-Json",
			collectionName, displayName, filePath,
		),
		"}",
		"Catch {",
		"$ErrorMessage = $_.Exception.Message;",
		"Write-Output -InputObject $ErrorMessage;",
		"exit 1;",
		"}",
	)

	if err != nil {
		return nil, err
	}

	out := res.Stdout

	return []byte(out), nil
}
