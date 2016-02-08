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

package apps

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
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

var GetAppsFailed = errors.New("Can't get apps list")
var UnpublishFailed = errors.New("Unpublish application failed")
var PublishFailed = errors.New("Publish application failed")
var AppsListUnavailable = errors.New("Apps list isn't available")
var FailedNameChange = errors.New("Failed to change the app name")

type Apps struct {
	user                 string
	server               string
	executionservers     []string
	sshport              string
	rdpport              string
	password             string
	windowsdomain        string
	xmlconfigurationfile string
	databaseuri          string
	protocol             string
	db                   *sql.DB
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

type Connection struct {
	Hostname  string `json:"hostname"`
	Port      string `json:"port"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	RemoteApp string `json:"remote_app"`
	Protocol  string `json:"protocol"`
	AppName   string `json:"app_name"`
}

func New(
	user,
	server,
	sshport,
	rdpport,
	protocol,
	password,
	databaseuri string,
	executionservers []string) *Apps {
	db := dbConnect(databaseuri)
	err := setupDb(db)
	if err != nil {
		panic(err)
	}

	go checkPublishedApps(user, password, sshport, server, db)
	return &Apps{
		user:             user,
		server:           server,
		sshport:          sshport,
		rdpport:          rdpport,
		protocol:         protocol,
		password:         password,
		databaseuri:      databaseuri,
		executionservers: executionservers,
		db:               db,
	}
}

func (a *Apps) GetAllApps() ([]ApplicationParams, error) {
	var applications []ApplicationParams
	rows, err := a.db.Query(
		`SELECT id, collection_name,
	alias, display_name,
	file_path
	FROM apps`,
	)

	if err != nil {
		log.Error("Connection to postgres failed: ", err.Error())
		return nil, GetAppsFailed
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
	return applications, nil

}

func (a *Apps) GetMyApps() ([]ApplicationParams, error) {
	var applications []ApplicationParams
	rows, err := a.db.Query(
		`SELECT id, collection_name,
	alias, display_name,
	file_path
	FROM apps`,
	)

	if err != nil {
		log.Error("Connection to postgres failed: ", err.Error())
		return nil, GetAppsFailed
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
		if appParam.Alias != "hapticPowershell" && appParam.DisplayName != "Desktop" {
			applications = append(applications, appParam)
		}
	}

	if len(applications) == 0 {
		applications = []ApplicationParams{}
	}
	return applications, nil

}

func dbConnect(dbURI string) *sql.DB {

	var err error
	var db *sql.DB

	for try := 0; try < 10; try++ {
		db, err = sql.Open("postgres", dbURI)
		if err == nil {
			return db
		}
		time.Sleep(time.Second * 5)
	}

	panic("Cannot connect to Postgres Database: " + err.Error())
}

// Connects to the postgres databse
func setupDb(db *sql.DB) error {
	rows, err := db.Query(
		`SELECT table_name
		FROM information_schema.tables
		WHERE table_name = 'apps'`)
	if err != nil {
		log.Error("Select tables names failed: ", err.Error())
		return err
	}
	defer rows.Close()

	if rows.Next() {
		log.Info("apps table already set up")
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
		log.Errorf("Unable to create apps table: %s", err)
		return err
	}

	rows.Close()
	return nil
}

func checkPublishedApps(user, password, sshport, server string, db *sql.DB) {
	for {
		time.Sleep(5 * time.Second)
		var applications []ApplicationParamsWin
		var apps []ApplicationParams
		cmd := exec.Command(
			"sshpass", "-p", password,
			"ssh", "-o", "StrictHostKeyChecking=no",
			"-p", sshport,
			fmt.Sprintf(
				"%s@%s",
				user,
				server,
			),
			"C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -Command \"Import-Module RemoteDesktop; Get-RDRemoteApp | ConvertTo-Json -Compress\"",
		)
		response, err := cmd.CombinedOutput()
		if err != nil {
			log.Error("Failed to execute sshpass command ", err, string(response))
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
					log.Error("Error inserting app into postgres: ", err.Error())
				}
			}
		}
		_, err = db.Query(
			`INSERT INTO apps
			(collection_name, alias, display_name, file_path)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar)
			`, "", "", "Desktop", "")
		if err != nil && !strings.Contains(err.Error(), "duplicate key") {
			log.Error("Error inserting hapticDesktop into postgres: ", err.Error())
		}

	}
}

// ========================================================================================================================
// Procedure: unpublishApplication
//
// Does:
// - Unpublish specified applications from ActiveDirectory
// ========================================================================================================================
func (a *Apps) UnpublishApp(Alias string) error {
	cmd := exec.Command(
		"sshpass", "-p", a.password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-p", a.sshport,
		fmt.Sprintf(
			"%s@%s",
			a.user,
			a.server,
		),
		"C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe -Command \"Import-Module RemoteDesktop; Remove-RDRemoteApp -Alias '"+Alias+"' -CollectionName collection -Force\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to execute sshpass command to unpublish an app", err, string(response))
		return UnpublishFailed
	}
	_, err = a.db.Query("DELETE FROM apps WHERE alias = $1::varchar", Alias)
	if err != nil {
		log.Error("delete from postgres failed: ", err)
		return UnpublishFailed
	}
	return nil
}

func (a *Apps) AppExists(alias string) (bool, error) {
	rows, err := a.db.Query(
		`SELECT alias
		FROM apps
		WHERE alias = $1::varchar`,
		alias)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (a *Apps) ChangeName(appId, newName string) error {
	_, err := a.db.Query(
		`UPDATE apps
		SET display_name = $1::varchar
		WHERE alias = $2::varchar`,
		newName, appId)
	if err != nil {
		log.Error("Changing app name failed: ", err)
		return FailedNameChange
	}
	return nil
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
func (a *Apps) PublishApp(path string) error {
	cmd := exec.Command(
		"sshpass", "-p", a.password,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-p", a.sshport,
		fmt.Sprintf(
			"%s@%s",
			a.user,
			a.server,
		),
		"C:/Windows/System32/WindowsPowerShell/v1.0/powershell.exe", "C:/publishApplication.ps1", path,
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to execute sshpass command to publish an app", err, string(response))
		return PublishFailed
	}
	return nil
}

func (a *Apps) RetrieveConnections(users []nano.User) ([]Connection, error) {

	// Seed random number generator
	rand.Seed(time.Now().UTC().UnixNano())
	var connections []Connection

	rows, err := a.db.Query(
		`SELECT alias
	FROM apps`,
	)

	if err != nil {
		log.Error("Unable to retrieve apps list from Postgres: ", err.Error())
		return nil, AppsListUnavailable
	}

	defer rows.Close()

	for rows.Next() {
		for _, user := range users {
			appParam := ApplicationParams{}

			rows.Scan(
				&appParam.Alias,
			)

			var conn Connection
			if count := len(a.executionservers); count > 0 {
				conn.Hostname = a.executionservers[rand.Intn(count)]
			} else {
				conn.Hostname = a.server
			}
			conn.Port = a.rdpport
			conn.Protocol = a.protocol
			if user.Sam != "" && appParam.Alias != "hapticPowershell" && appParam.Alias != "" {
				conn.Username = user.Sam
				conn.Password = user.WindowsPassword
			} else {
				continue
			}
			conn.RemoteApp = "||" + appParam.Alias
			conn.AppName = appParam.Alias
			connections = append(connections, conn)
		}
	}

	connections = append(connections, Connection{
		Hostname:  a.server,
		Port:      a.rdpport,
		Protocol:  a.protocol,
		Username:  a.user,
		Password:  a.password,
		RemoteApp: "",
		AppName:   "hapticDesktop",
	})
	connections = append(connections, Connection{
		Hostname:  a.server,
		Port:      a.rdpport,
		Protocol:  a.protocol,
		Username:  a.user,
		Password:  a.password,
		RemoteApp: "||hapticPowershell",
		AppName:   "hapticPowershell",
	})

	return connections, nil
}
