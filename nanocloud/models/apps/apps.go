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
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
)

var (
	GetAppsFailed       = errors.New("Can't get apps list")
	UnpublishFailed     = errors.New("Unpublish application failed")
	PublishFailed       = errors.New("Publish application failed")
	AppsListUnavailable = errors.New("Apps list isn't available")
	FailedNameChange    = errors.New("Failed to change the app name")
)

var (
	kUser                 string
	kServer               string
	kExecutionServers     []string
	kSSHPort              string
	kRDPPort              string
	kPassword             string
	kWindowsDomain        string
	kXMLConfigurationFile string
	kProtocol             string
)

type ApplicationParams struct {
	Id             int    `json:"-"`
	CollectionName string `json:"collection-name"`
	Alias          string `json:"alias"`
	DisplayName    string `json:"display-name"`
	FilePath       string `json:"file-path"`
	IconContents   []byte `json:"icon-content"`
}

type ApplicationParamsWin struct {
	Id             int
	CollectionName string
	Alias          string
	DisplayName    string
	FilePath       string
	IconContents   []byte
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

func AppExists(appId string) (bool, error) {
	rows, err := db.Query(
		`SELECT alias
     FROM apps
     WHERE id = $1::int`,
		appId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func ChangeName(appId, newName string) error {
	_, err := db.Query(
		`UPDATE apps
     SET display_name = $1::varchar
     WHERE id = $2::int`,
		newName, appId)
	if err != nil {
		log.Error("Changing app name failed: ", err)
		return FailedNameChange
	}
	return nil
}

func GetAllApps() ([]ApplicationParams, error) {
	var applications []ApplicationParams
	rows, err := db.Query(
		`SELECT id, collection_name,
	alias, display_name,
	file_path,
	icon_content
	FROM apps`)

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
			&appParam.IconContents,
		)
		applications = append(applications, appParam)

	}

	if len(applications) == 0 {
		applications = []ApplicationParams{}
	}
	return applications, nil
}

func GetUserApps(userId string) ([]ApplicationParams, error) {
	var applications []ApplicationParams
	rows, err := db.Query(
		`SELECT id, collection_name,
	alias, display_name,
	file_path,
	icon_content
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
			&appParam.IconContents,
		)
		if appParam.Alias != "hapticPowershell" && appParam.Alias != "Desktop" {
			applications = append(applications, appParam)
		}
	}

	if len(applications) == 0 {
		applications = []ApplicationParams{}
	}
	return applications, nil
}

func CheckPublishedApps() {
	_, err := db.Query(
		`INSERT INTO apps
			(collection_name, alias, display_name, file_path, icon_content)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::bytea)
			`, "", "Desktop", "Desktop", "", "")
	if err != nil && !strings.Contains(err.Error(), "duplicate key") {
		log.Error("Error inserting hapticDesktop into postgres: ", err.Error())
	}
	for {
		time.Sleep(5 * time.Second)
		var applications []ApplicationParamsWin
		var winapp ApplicationParamsWin
		var apps []ApplicationParams
		cmd := exec.Command(
			"sshpass", "-p", kPassword,
			"ssh", "-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
			"-o", "LogLevel=quiet",
			"-p", kSSHPort,
			fmt.Sprintf(
				"%s@%s",
				kUser,
				kServer,
			),
			"powershell.exe \"Import-Module RemoteDesktop; Get-RDRemoteApp | ConvertTo-Json -Compress\"",
		)
		response, err := cmd.CombinedOutput()

		if err != nil {
			log.Error("Failed to execute sshpass command ", err, string(response))
			continue
		}
		err = json.Unmarshal(response, &applications)
		if err != nil {

			err = json.Unmarshal(response, &winapp)
			if err != nil {
				continue
			}
			application := ApplicationParams{
				CollectionName: winapp.CollectionName,
				DisplayName:    winapp.DisplayName,
				Alias:          winapp.Alias,
				FilePath:       winapp.FilePath,
				IconContents:   winapp.IconContents,
			}

			_, err := db.Query(
				`INSERT INTO apps
			(collection_name, alias, display_name, file_path, icon_content)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::bytea)
			`, application.CollectionName, application.Alias, application.DisplayName, application.FilePath, application.IconContents)
			if err != nil && !strings.Contains(err.Error(), "duplicate key") {
				log.Error("Error inserting app into postgres: ", err.Error())
			}
			continue
		}
		for _, app := range applications {
			apps = append(apps, ApplicationParams{
				CollectionName: app.CollectionName,
				DisplayName:    app.DisplayName,
				Alias:          app.Alias,
				FilePath:       app.FilePath,
				IconContents:   app.IconContents,
			})
		}

		for _, application := range apps {

			if application.CollectionName != "" && application.Alias != "" && application.DisplayName != "" && application.FilePath != "" {
				_, err := db.Query(
					`INSERT INTO apps
			(collection_name, alias, display_name, file_path, icon_content)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::bytea)
			`, application.CollectionName, application.Alias, application.DisplayName, application.FilePath, application.IconContents)
				if err != nil && !strings.Contains(err.Error(), "duplicate key") {
					log.Error("Error inserting app into postgres: ", err.Error())
				}
			}
		}
	}
}

// ========================================================================================================================
// Procedure: unpublishApplication
//
// Does:
// - Unpublish specified applications from ActiveDirectory
// ========================================================================================================================
func UnpublishApp(appId string) error {

	res, err := db.Query("SELECT alias FROM apps WHERE id = $1::int", appId)
	if err != nil {
		return err
	}
	defer res.Close()
	var Alias string
	for res.Next() {
		err := res.Scan(&Alias)
		if err != nil {
			continue
		}
	}

	cmd := exec.Command(
		"sshpass", "-p", kPassword,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", kSSHPort,
		fmt.Sprintf(
			"%s@%s",
			kUser,
			kServer,
		),
		"powershell.exe \"Import-Module RemoteDesktop; Remove-RDRemoteApp -Alias '"+Alias+"' -CollectionName collection -Force\"",
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to execute sshpass command to unpublish an app", err, string(response))
		return UnpublishFailed
	}
	_, err = db.Query("DELETE FROM apps WHERE alias = $1::varchar", Alias)
	if err != nil {
		log.Error("delete from postgres failed: ", err)
		return UnpublishFailed
	}
	return nil
}

func PublishApp(path string) error {
	cmd := exec.Command(
		"sshpass", "-p", kPassword,
		"ssh", "-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-p", kSSHPort,
		fmt.Sprintf(
			"%s@%s",
			kUser,
			kServer,
		),
		"powershell.exe -file C:\\publishApplication.ps1 "+path,
	)
	response, err := cmd.CombinedOutput()
	if err != nil {
		log.Error("Failed to execute sshpass command to publish an app", err, string(response))
		return PublishFailed
	}
	return nil
}

func RetrieveConnections(user *users.User, users *[]users.User) ([]Connection, error) {

	rand.Seed(time.Now().UTC().UnixNano())
	var connections []Connection

	rows, err := db.Query("SELECT alias FROM apps")
	if err != nil {
		log.Error("Unable to retrieve apps list from Postgres: ", err.Error())
		return nil, AppsListUnavailable
	}
	defer rows.Close()
	var execServ string
	for rows.Next() {
		appParam := ApplicationParams{}
		rows.Scan(
			&appParam.Alias,
		)

		if count := len(kExecutionServers); count > 0 {
			execServ = kExecutionServers[rand.Intn(count)]
		} else {
			execServ = kServer
		}
		username := user.Sam + "@" + kWindowsDomain
		pwd := user.WindowsPassword
		var conn Connection
		if appParam.Alias != "Desktop" {
			conn = Connection{
				Hostname:  execServ,
				Port:      kRDPPort,
				Protocol:  kProtocol,
				Username:  username,
				Password:  pwd,
				RemoteApp: "||" + appParam.Alias,
				AppName:   appParam.Alias,
			}
		} else {
			conn = Connection{
				Hostname:  execServ,
				Port:      kRDPPort,
				Protocol:  kProtocol,
				Username:  username,
				Password:  pwd,
				RemoteApp: "",
				AppName:   "hapticDesktop",
			}
		}
		connections = append(connections, conn)
	}

	return connections, nil
}

func init() {
	kUser = utils.Env("USER", "Administrator")
	kProtocol = utils.Env("PROTOCOL", "rdp")
	kSSHPort = utils.Env("SSH_PORT", "22")
	kRDPPort = utils.Env("RDP_PORT", "3389")
	kServer = utils.Env("SERVER", "62.210.56.76")
	kPassword = utils.Env("PASSWORD", "ItsPass1942+")
	kWindowsDomain = utils.Env("WINDOWS_DOMAIN", "intra.localdomain.com")
	kExecutionServers = strings.Split(utils.Env("EXECUTION_SERVERS", "62.210.56.76"), ",")
}
