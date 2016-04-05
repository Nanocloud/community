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
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
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

func AddApp(params ApplicationParams) error {
	_, err := db.Query(
		`INSERT INTO apps
			(collection_name, alias, display_name, file_path)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar)
			`, params.CollectionName, params.Alias, params.DisplayName, params.FilePath)
	if err != nil && !strings.Contains(err.Error(), "duplicate key") {
		log.Error("Error inserting app into postgres: ", err.Error())
	}
	return nil
}

func CheckPublishedApps() {
	_, err := db.Query(
		`INSERT INTO apps
			(collection_name, alias, display_name, file_path, icon_content)
			VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::bytea)
			`, "", "hapticDesktop", "Desktop", "", "")
	if err != nil && !strings.Contains(err.Error(), "duplicate key") {
		log.Error("Error inserting hapticDesktop into postgres: ", err.Error())
	}
	for {
		time.Sleep(5 * time.Second)
		var winapp ApplicationParams
		var apps []ApplicationParams
		resp, err := http.Get("http://" + kServer + ":9090/apps")
		if err != nil {
			continue
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			continue
		}
		err = json.Unmarshal(b, &apps)
		if err != nil {
			err = json.Unmarshal(b, &winapp)
			if err != nil {
				continue
			}
			if winapp.CollectionName != "" && winapp.Alias != "" && winapp.DisplayName != "" && winapp.FilePath != "" {
				_, err := db.Query(
					`INSERT INTO apps
				(collection_name, alias, display_name, file_path, icon_content)
				VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::bytea)
				`, winapp.CollectionName, winapp.Alias, winapp.DisplayName, winapp.FilePath, winapp.IconContents)
				if err != nil && !strings.Contains(err.Error(), "duplicate key") {
					log.Error("Error inserting app into postgres: ", err.Error())
				}
			}
			continue
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

func getCredentials() (string, string) {
	rows, err := db.Query(
		`SELECT sam, windows_password
		FROM users
		WHERE is_admin = $1::boolean`,
		true)
	if err != nil {
		return "", ""
	}
	var name string
	var pwd string
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&name, &pwd)
		if err != nil {
			return "", ""
		}
	}
	return name, pwd
}

// ========================================================================================================================
// Procedure: unpublishApplication
//
// Does:
// - Unpublish specified applications from ActiveDirectory
// ========================================================================================================================
func UnpublishApp(Alias string) error {
	id, err := strconv.Atoi(Alias)
	if err != nil {
		return UnpublishFailed
	}
	rows, err := db.Query(
		`SELECT alias FROM apps WHERE id = $1::int`,
		id,
	)
	if err != nil {
		log.Error("Connection to postgres failed: ", err.Error())
		return UnpublishFailed
	}

	defer rows.Close()

	var alias string
	for rows.Next() {
		rows.Scan(
			&alias,
		)
	}
	log.Error(alias)

	req, err := http.NewRequest("DELETE", "http://"+kServer+":9090/apps/"+alias, nil)
	username, pwd := getCredentials()
	if username == "" || pwd == "" {
		log.Error("Unable to retrieve admin credentials")
		return UnpublishFailed
	}
	req.SetBasicAuth(username, pwd)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return UnpublishFailed
	}
	if resp.Status != "200 OK" {
		log.Error("Plaza return code: " + resp.Status)
		return UnpublishFailed
	}
	_, err = db.Query("DELETE FROM apps WHERE alias = $1::varchar", alias)
	if err != nil {
		log.Error("delete from postgres failed: ", err)
		return UnpublishFailed
	}
	return nil
}

func PublishApp(body io.Reader) error {
	req, err := http.NewRequest("POST", "http://"+kServer+":9090/publishapp", body)
	username, pwd := getCredentials()
	if username == "" || pwd == "" {
		log.Error("Unable to retrieve admin credentials")
		return PublishFailed
	}
	req.SetBasicAuth(username, pwd)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return PublishFailed
	}
	if resp.Status != "200 OK" {
		log.Error("Plaza return code: " + resp.Status)
		return PublishFailed
	}
	return nil
}

func RetrieveConnections(user *users.User, users []*users.User) ([]Connection, error) {

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
		if appParam.Alias != "hapticDesktop" {
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
