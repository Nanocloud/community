/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
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

package main

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"time"

	"github.com/Nanocloud/nano"
)

const (
	windowsUserPassword = "12345abcDEF+"
)

type GuacamoleXMLConfigs struct {
	XMLName xml.Name             `xml:configs`
	Config  []GuacamoleXMLConfig `xml:"config"`
}

type GuacamoleXMLConfig struct {
	XMLName  xml.Name            `xml:config`
	Name     string              `xml:"name,attr"`
	Protocol string              `xml:"protocol,attr"`
	Params   []GuacamoleXMLParam `xml:"param"`
}

type GuacamoleXMLParam struct {
	ParamName  string `xml:"name,attr"`
	ParamValue string `xml:"value,attr"`
}

type Connection struct {
	Hostname       string `xml:"hostname"`
	Port           string `xml:"port"`
	Username       string `xml:"username"`
	Password       string `xml:"password"`
	RemoteApp      string `xml:"remote-app"`
	ConnectionName string
}

type ApplicationParams struct {
	CollectionName string
	Alias          string
	DisplayName    string
	IconContents   []uint8
	FilePath       string
}

func getUsers() ([]nano.User, error) {
	res, err := module.Request("GET", "/users", "", nil, nil)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("invalid status code")
	}

	var m []nano.User
	module.Log.Error(string(res.Body))
	err = json.Unmarshal(res.Body, &m)

	if err != nil {
		module.Log.Error(err)
		return nil, err
	}

	return m, nil
}

// ========================================================================================================================
// Procedure: createConnections
//
// Does:
// - Create all connections in DB for a particular user in order to use all applications
// ========================================================================================================================
func createConnections() error {

	type configs GuacamoleXMLConfigs
	var (
		applications    []ApplicationParams
		connections     configs
		executionServer string
	)

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	cmd := exec.Command(
		"sshpass", "-p", conf.Password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-p", conf.SSHPort,
		fmt.Sprintf(
			"%s@%s",
			conf.User,
			conf.Server,
		),
		"C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -Command \"Import-Module RemoteDesktop; Get-RDRemoteApp | ConvertTo-Json -Compress\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		module.Log.Error("Failed to execute sshpass command ", err, string(response))
		return err
	}

	if []byte(response)[0] != []byte("[")[0] {
		response = []byte(fmt.Sprintf("[%s]", string(response)))
	}
	json.Unmarshal(response, &applications)
	for _, application := range applications {
		application.IconContents = []byte(base64.StdEncoding.EncodeToString(application.IconContents))
	}

	//	users, _ := g_Db.GetUsers()
	users, err := getUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		for _, application := range applications {
			if application.Alias == "hapticPowershell" {
				continue
			}

			// Select randomly execution machine from availbale execution machines
			if count := len(conf.ExecutionServers); count > 0 {
				executionServer = conf.ExecutionServers[rand.Intn(count)]
			} else {
				executionServer = conf.Server
			}

			connections.Config = append(connections.Config, GuacamoleXMLConfig{
				Name:     fmt.Sprintf("%s_%s", application.Alias, user.Email),
				Protocol: "rdp",
				Params: []GuacamoleXMLParam{
					GuacamoleXMLParam{
						ParamName:  "hostname",
						ParamValue: executionServer,
					},
					GuacamoleXMLParam{
						ParamName:  "port",
						ParamValue: conf.RDPPort,
					},
					GuacamoleXMLParam{
						ParamName:  "username",
						ParamValue: user.Sam,
					},
					GuacamoleXMLParam{
						ParamName:  "password",
						ParamValue: user.WindowsPassword,
					},
					GuacamoleXMLParam{
						ParamName:  "remote-app",
						ParamValue: fmt.Sprintf("||%s", application.Alias),
					},
					GuacamoleXMLParam{
						ParamName:  "security",
						ParamValue: "nla",
					},
					GuacamoleXMLParam{
						ParamName:  "ignore-cert",
						ParamValue: "true",
					},
				},
			})
		}
	}

	connections.Config = append(connections.Config, GuacamoleXMLConfig{
		Name:     "hapticDesktop",
		Protocol: "rdp",
		Params: []GuacamoleXMLParam{
			GuacamoleXMLParam{
				ParamName:  "hostname",
				ParamValue: conf.Server,
			},
			GuacamoleXMLParam{
				ParamName:  "port",
				ParamValue: conf.RDPPort,
			},
			GuacamoleXMLParam{
				ParamName:  "username",
				ParamValue: conf.User,
			},
			GuacamoleXMLParam{
				ParamName:  "password",
				ParamValue: conf.Password,
			},
			GuacamoleXMLParam{
				ParamName:  "security",
				ParamValue: "nla",
			},
			GuacamoleXMLParam{
				ParamName:  "ignore-cert",
				ParamValue: "true",
			},
		},
	})
	connections.Config = append(connections.Config, GuacamoleXMLConfig{
		Name:     "hapticPowershell",
		Protocol: "rdp",
		Params: []GuacamoleXMLParam{
			GuacamoleXMLParam{
				ParamName:  "hostname",
				ParamValue: conf.Server,
			},
			GuacamoleXMLParam{
				ParamName:  "port",
				ParamValue: conf.RDPPort,
			},
			GuacamoleXMLParam{
				ParamName:  "username",
				ParamValue: conf.User,
			},
			GuacamoleXMLParam{
				ParamName:  "password",
				ParamValue: conf.Password,
			},
			GuacamoleXMLParam{
				ParamName:  "remote-app",
				ParamValue: "||hapticPowershell",
			},
			GuacamoleXMLParam{
				ParamName:  "security",
				ParamValue: "nla",
			},
			GuacamoleXMLParam{
				ParamName:  "ignore-cert",
				ParamValue: "true",
			},
		},
	})

	output, err := xml.MarshalIndent(connections, "  ", "    ")
	if err != nil {
		module.Log.Error("xml Marshalling of connections failed: ", err)
		return err
	}

	if err = ioutil.WriteFile(conf.XMLConfigurationFile, output, 0777); err != nil {
		module.Log.Error("Failed to save connections in ", conf.XMLConfigurationFile, " params: ", err)
		return err
	}

	return nil
}

// ========================================================================================================================
// Procedure: listApplications
//
// Does:
// - Return list of applications published by Active Directory
// ========================================================================================================================
func listApplications(req nano.Request) (*nano.Response, error) {
	var (
		guacamoleConfigs GuacamoleXMLConfigs
		connections      []Connection
		bytesRead        []byte
		err              error
	)

	err = createConnections()
	if err != nil {
		return nil, err
	}

	if bytesRead, err = ioutil.ReadFile(conf.XMLConfigurationFile); err != nil {
		module.Log.Error("Failed to read connections params in XMLConfigurationFile: ", err)
		return nil, err
	}

	err = xml.Unmarshal(bytesRead, &guacamoleConfigs)
	if err != nil {
		return nil, err
	}

	for _, config := range guacamoleConfigs.Config {
		var connection Connection

		for _, param := range config.Params {
			switch true {
			case param.ParamName == "hostname":
				connection.Hostname = param.ParamValue
			case param.ParamName == "port":
				connection.Port = param.ParamValue
			case param.ParamName == "username":
				connection.Username = param.ParamValue
			case param.ParamName == "password":
				connection.Password = param.ParamValue
			case param.ParamName == "remote-app":
				connection.RemoteApp = param.ParamValue
			}
		}
		connection.ConnectionName = config.Name

		if connection.RemoteApp == "" || connection.RemoteApp == "||hapticPowershell" {
			continue
		}

		connections = append(connections, connection)
	}
	return nano.JSONResponse(200, connections), nil
}

// ========================================================================================================================
// Procedure: listApplicationsForSamAccount
//
// Does:
// - Return list of applications available for a particular SAM account
// ========================================================================================================================
func listApplicationsForSamAccount(req nano.Request) (*nano.Response, error) {
	var (
		guacamoleConfigs GuacamoleXMLConfigs
		connections      []Connection
		bytesRead        []byte
		err              error
	)

	if bytesRead, err = ioutil.ReadFile(conf.XMLConfigurationFile); err != nil {
		module.Log.Error("Failed to read connections params in XMLConfigurationFile: ", err)
		return nil, err
	}

	err = xml.Unmarshal(bytesRead, &guacamoleConfigs)
	if err != nil {
		return nil, err
	}

	for _, config := range guacamoleConfigs.Config {
		var connection Connection

		if connection.ConnectionName == "hapticPowershell" {
			continue
		}

		connection.ConnectionName = config.Name
		for _, param := range config.Params {
			switch true {
			case param.ParamName == "hostname":
				connection.Hostname = param.ParamValue
			case param.ParamName == "port":
				connection.Port = param.ParamValue
			case param.ParamName == "username":
				connection.Username = param.ParamValue
			case param.ParamName == "password":
				connection.Password = param.ParamValue
			case param.ParamName == "remote-app":
				connection.RemoteApp = param.ParamValue
			}
		}

		if connection.Username == fmt.Sprintf("%s@%s", req.User.Sam, conf.WindowsDomain) {
			connections = append(connections, connection)
		}
	}

	return nano.JSONResponse(200, connections), nil
}

// ========================================================================================================================
// Procedure: unpublishApplication
//
// Does:
// - Unpublish specified applications from ActiveDirectory
// ========================================================================================================================
func unpublishApp(Alias string) error {
	cmd := exec.Command(
		"sshpass", "-p", conf.Password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-p", conf.SSHPort,
		fmt.Sprintf(
			"%s@%s",
			conf.User,
			conf.Server,
		),
		"C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -Command \"Import-Module RemoteDesktop; Remove-RDRemoteApp -Alias "+Alias+" -CollectionName collection -Force\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		module.Log.Error("Failed to execute sshpass command to unpublish an app", err, string(response))
	}
	return err
}

/*
// ========================================================================================================================
// Procedure: SyncUploadedFile
//
// Does:
// - Upload user files to windows VM
// ========================================================================================================================
func syncUploadedFile(Filename string) {
	bashCopyScript := filepath.Join(nan.Config().CommonBaseDir, "scripts", "copy.sh")
	cmd := exec.Command(bashCopyScript, Filename)
	cmd.Dir = filepath.Join(nan.Config().CommonBaseDir, "scripts")
	response, err := cmd.Output()
	if err != nil {
		LogError("Failed to run script copy.sh, error: %s, output: %s\n", err, string(response))
	} else {
		Log("SCP upload success for file %s\n", Filename)
	}
}*/
