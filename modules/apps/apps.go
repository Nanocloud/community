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
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/Nanocloud/nano"
)

var db *sql.DB

type Connection struct {
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	RemoteApp string `json:"remote_app"`
	Protocol  string `json:"protocol"`
}

type ApplicationParams struct {
	Id             int    `json:"id"`
	CollectionName string `json:"collection_name"`
	Alias          string `json:"alias"`
	DisplayName    string `json:"display_name"`
	FilePath       string `json:"file_path"`
}

type ApplicationParamsWin struct {
	Id             int
	CollectionName string
	Alias          string
	DisplayName    string
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
func getConnections(req nano.Request) (*nano.Response, error) {

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())
	var connections []Connection

	rows, err := db.Query(
		`SELECT alias
	FROM apps`,
	)

	if err != nil {
		module.Log.Error(err.Error())
		return nil, err
	}

	defer rows.Close()

	users, err := getUsers()
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": "Unable to retrieve users" + err.Error(),
		}), nil
	}

	for rows.Next() {
		for _, user := range users {
			module.Log.Error(user)
			appParam := ApplicationParams{}

			rows.Scan(
				&appParam.Alias,
			)

			var conn Connection
			if count := len(conf.ExecutionServers); count > 0 {
				conn.Hostname = conf.ExecutionServers[rand.Intn(count)]
			} else {
				conn.Hostname = conf.Server
			}
			conn.Port = conf.RDPPort
			conn.Protocol = conf.Protocol
			if user.Sam != "" && appParam.Alias != "hapticPowershell" && appParam.Alias != "" {
				conn.Username = user.Sam
				conn.Password = user.WindowsPassword
			} else {
				continue
			}
			conn.RemoteApp = "||" + appParam.Alias
			connections = append(connections, conn)
		}
	}

	connections = append(connections, Connection{
		Hostname:  conf.Server,
		Port:      conf.RDPPort,
		Protocol:  conf.Protocol,
		Username:  conf.User,
		Password:  conf.Password,
		RemoteApp: "",
	})
	connections = append(connections, Connection{
		Hostname:  conf.Server,
		Port:      conf.RDPPort,
		Protocol:  conf.Protocol,
		Username:  conf.User,
		Password:  conf.Password,
		RemoteApp: "||hapticPowershell",
	})
	return nano.JSONResponse(200, connections), nil
}

func listApplications(req nano.Request) (*nano.Response, error) {

	var applications []ApplicationParams
	rows, err := db.Query(
		`SELECT id, collection_name,
	alias, display_name,
	file_path
	FROM apps`,
	)

	if err != nil {
		module.Log.Error("Connection to postgres failed: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		appParam := ApplicationParams{}

		rows.Scan(
			&appParam.Id,
			&appParam.CollectionName,
			&appParam.Alias,
			&appParam.DisplayName,
			&appParam.FilePath,
		)
		applications = append(applications, appParam)

	}

	if len(applications) == 0 {
		applications = []ApplicationParams{}
	}

	return nano.JSONResponse(200, applications), nil
}

// ========================================================================================================================
// Procedure: listApplicationsForSamAccount
//
// Does:
// - Return list of applications available for a particular SAM account
// ========================================================================================================================
func listApplicationsForSamAccount(req nano.Request) (*nano.Response, error) {

	var applications []ApplicationParams
	rows, err := db.Query(
		`SELECT id, collection_name,
	alias, display_name,
	file_path
	FROM apps`,
	)

	if err != nil {
		module.Log.Error("Connection to postgres failed: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		appParam := ApplicationParams{}

		rows.Scan(
			&appParam.Id,
			&appParam.CollectionName,
			&appParam.Alias,
			&appParam.DisplayName,
			&appParam.FilePath,
		)

		//TODO : ONLY APPEND IF USER GROUP HAS ACCES TO THE APP
		if appParam.Alias != "hapticPowershell" && appParam.DisplayName != "Desktop" {
			applications = append(applications, appParam)
		}

	}

	if len(applications) == 0 {
		applications = []ApplicationParams{}
	}

	return nano.JSONResponse(200, applications), nil

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
		"C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -Command \"Import-Module RemoteDesktop; Remove-RDRemoteApp -Alias '"+Alias+"' -CollectionName collection -Force\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		module.Log.Error("Failed to execute sshpass command to unpublish an app", err, string(response))
		return err
	}
	_, err = db.Query("DELETE FROM apps WHERE alias = $1::varchar", Alias)
	if err != nil {
		module.Log.Error("delete from postgres failed: ", err)
		return err
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

func checkPublishedApps() {
	for {
		time.Sleep(5 * time.Second)
		var applications []ApplicationParamsWin
		var apps []ApplicationParams
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
			continue
		}

		err = json.Unmarshal(response, &applications)
		if err != nil {
			continue
		}
		for _, app := range applications {
			apps = append(apps, ApplicationParams{
				CollectionName: app.CollectionName,
				DisplayName:    app.DisplayName,
				Alias:          app.Alias,
				FilePath:       app.FilePath,
			})
		}

		for _, application := range apps {
			if application.CollectionName != "" && application.Alias != "" && application.DisplayName != "" && application.FilePath != "" {
				_, err := db.Query(
					`INSERT INTO apps
			(collection_name, alias, display_name, file_path)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar)
			`, application.CollectionName, application.Alias, application.DisplayName, application.FilePath)
				if err != nil && !strings.Contains(err.Error(), "duplicate key") {
					module.Log.Error("Error inserting app into postgres: ", err.Error())
				}
			}
		}
		_, err = db.Query(
			`INSERT INTO apps
			(collection_name, alias, display_name, file_path)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar)
			`, "", "", "Desktop", "")
		if err != nil && !strings.Contains(err.Error(), "duplicate key") {
			module.Log.Error("Error inserting hapticDesktop into postgres: ", err.Error())
		}

	}
}

func publishApp(path string) error {
	cmd := exec.Command(
		"sshpass", "-p", conf.Password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-p", conf.SSHPort,
		fmt.Sprintf(
			"%s@%s",
			conf.User,
			conf.Server,
		),
		"C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe", "C:/publishApplication.ps1", path,
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		module.Log.Error("Failed to execute sshpass command to publish an app", err, string(response))
	}
	return err
}

func dbConnect() {

	var err error

	for try := 0; try < 10; try++ {
		db, err = sql.Open("postgres", conf.DatabaseURI)
		if err == nil {
			return
		}
		time.Sleep(time.Second * 5)
	}

	module.Log.Fatalf("Cannot connect to Postgres Database: %s", err)
}

// Connects to the postgres databse
func setupDb() error {
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'apps'`)
	if err != nil {
		module.Log.Error(err.Error())
		return err
	}
	defer rows.Close()

	if rows.Next() {
		module.Log.Info("apps table already set up")
		return nil
	}
	rows, err = db.Query(
		`CREATE TABLE apps (
			id	SERIAL PRIMARY KEY,
			collection_name		varchar(36),
			alias		varchar(36) UNIQUE,
			display_name		varchar(36),
			file_path		 varchar(255)
		);`)
	if err != nil {
		module.Log.Errorf("Unable to create apps table: %s", err)
		return err
	}

	rows.Close()
	return nil
}
