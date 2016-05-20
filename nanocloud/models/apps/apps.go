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
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/plaza"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
	uuid "github.com/satori/go.uuid"
)

var (
	GetAppsFailed       = errors.New("Can't get apps list")
	GetAppFailed        = errors.New("Can't get app")
	UnpublishFailed     = errors.New("Unpublish application failed")
	PublishFailed       = errors.New("Publish application failed")
	AppsListUnavailable = errors.New("Apps list isn't available")
	FailedNameChange    = errors.New("Failed to change the app name")
)

var (
	kServer               string
	kExecutionServers     []string
	kRDPPort              string
	kXMLConfigurationFile string
	kProtocol             string
)

type Application struct {
	Id             string `json:"-"`
	CollectionName string `json:"collection-name"`
	Alias          string `json:"alias"`
	DisplayName    string `json:"display-name"`
	FilePath       string `json:"file-path"`
	Path           string `json:"path"`
	IconContents   []byte `json:"icon-content"`
}

func (a *Application) GetID() string {
	return a.Id
}

func (a *Application) SetID(id string) error {
	a.Id = id
	return nil
}

type ApplicationWin struct {
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
		WHERE id = $1::varchar`,
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
		WHERE id = $2::varchar`,
		newName, appId)
	if err != nil {
		log.Error("Changing app name failed: ", err)
		return FailedNameChange
	}
	return nil
}

func GetApp(appId string) (*Application, error) {
	rows, err := db.Query(
		`SELECT id, collection_name,
		alias, display_name,
		file_path,
		icon_content
		FROM apps WHERE id = $1::varchar`, appId)

	if err != nil {
		return nil, GetAppFailed
	}

	defer rows.Close()
	if rows.Next() {
		var application Application

		err = rows.Scan(
			&application.Id,
			&application.CollectionName,
			&application.Alias,
			&application.DisplayName,
			&application.FilePath,
			&application.IconContents,
		)
		if err != nil {
			return nil, err
		}

		return &application, nil
	}
	return nil, nil
}

func GetAllApps() ([]*Application, error) {
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

	applications := make([]*Application, 0)

	for rows.Next() {
		appParam := Application{}

		rows.Scan(
			&appParam.Id,
			&appParam.CollectionName,
			&appParam.Alias,
			&appParam.DisplayName,
			&appParam.FilePath,
			&appParam.IconContents,
		)
		applications = append(applications, &appParam)

	}

	return applications, nil

}

func GetUserApps(userId string) ([]*Application, error) {
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

	applications := make([]*Application, 0)

	for rows.Next() {
		appParam := Application{}

		rows.Scan(
			&appParam.Id,
			&appParam.CollectionName,
			&appParam.Alias,
			&appParam.DisplayName,
			&appParam.FilePath,
			&appParam.IconContents,
		)
		if appParam.Alias != "hapticPowershell" && appParam.Alias != "Desktop" {
			applications = append(applications, &appParam)
		}
	}

	return applications, nil
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
func UnpublishApp(user *users.User, id string) error {
	rows, err := db.Query(
		`SELECT alias, collection_name FROM apps WHERE id = $1::varchar`,
		id,
	)
	if err != nil {
		log.Error(err)
		return UnpublishFailed
	}

	defer rows.Close()

	var alias string
	var collection string
	if !rows.Next() {
		return errors.New("Application not found")
	}

	rows.Scan(&alias, &collection)

	plazaAddress := utils.Env("EXECUTION_SERVERS", "iaas-module")
	if plazaAddress == "" {
		return errors.New("plaza address unknown")
	}

	plazaPort, err := strconv.Atoi(utils.Env("PLAZA_PORT", "9090"))
	if err != nil {
		return err
	}

	winUser, err := user.WindowsCredentials()
	if err != nil {
		return err
	}

	_, err = plaza.UnpublishApp(
		plazaAddress, plazaPort,
		winUser.Sam,
		winUser.Domain,
		winUser.Password,
		collection,
		alias,
	)

	if err != nil {
		log.Error(err)
		return UnpublishFailed
	}

	_, err = db.Query("DELETE FROM apps WHERE id = $1::varchar", id)
	if err != nil {
		log.Error("delete from postgres failed: ", err)
		return UnpublishFailed
	}
	return nil
}

func PublishApp(user *users.User, app *Application) error {
	plazaAddress := utils.Env("EXECUTION_SERVERS", "iaas-module")
	if plazaAddress == "" {
		return errors.New("plaza address unknown")
	}

	plazaPort, err := strconv.Atoi(utils.Env("PLAZA_PORT", "9090"))
	if err != nil {
		return err
	}

	winUser, err := user.WindowsCredentials()
	if err != nil {
		return err
	}

	res, err := plaza.PublishApp(
		plazaAddress, plazaPort,
		winUser.Sam,
		winUser.Domain,
		winUser.Password,
		app.CollectionName,
		app.DisplayName,
		app.Path,
	)

	if err != nil {
		log.Error(err)
		return PublishFailed
	}

	a := ApplicationWin{}
	err = json.Unmarshal(res, &a)
	if err != nil {
		return err
	}

	id := uuid.NewV4().String()

	_, err = db.Query(
		`INSERT INTO apps
		(id, collection_name, alias, display_name, file_path, icon_content)
		VALUES ( $1::varchar, $2::varchar, $3::varchar, $4::varchar, $5::varchar, $6::bytea)
		`,
		id, a.CollectionName, a.Alias, a.DisplayName, a.FilePath, a.IconContents,
	)

	if err != nil {
		return err
	}
	app.CollectionName = a.CollectionName
	app.Alias = a.Alias
	app.DisplayName = a.DisplayName
	app.FilePath = a.FilePath
	app.IconContents = a.IconContents
	app.Id = id

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
		appParam := Application{}
		rows.Scan(
			&appParam.Alias,
		)

		if count := len(kExecutionServers); count > 0 {
			execServ = kExecutionServers[rand.Intn(count)]
		} else {
			execServ = kServer
		}

		winUser, err := user.WindowsCredentials()
		if err != nil {
			return nil, err
		}

		username := winUser.Sam + "@" + winUser.Domain
		pwd := winUser.Password

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
	kProtocol = utils.Env("PROTOCOL", "rdp")
	kRDPPort = utils.Env("RDP_PORT", "3389")
	kServer = utils.Env("EXECUTION_SERVERS", "iaas-module")
	kExecutionServers = strings.Split(utils.Env("EXECUTION_SERVERS", ""), ",")

	if kServer == "" {
		panic("EXECUTION_SERVERS not set")
	}

	if len(kExecutionServers) == 0 {
		panic("EXECUTION_SERVERS not set")
	}
}
