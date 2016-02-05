/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"encoding/json"
	"net/url"
	"os"

	ldapPkg "github.com/Nanocloud/community/modules/ldap/lib/ldap"
	"github.com/Nanocloud/nano"
)

var module nano.Module

var conf struct {
	Username   string
	Password   string
	ServerURL  string
	Ou         string
	LDAPServer url.URL
}

type handler struct {
	ldapCon *ldapPkg.Ldap
}

type hash map[string]interface{}

type AccountParams struct {
	UserEmail string
	Password  string
}

type ChangePasswordParams struct {
	SamAccountName string
	NewPassword    string
}

// Strucutre used in messages from RabbitMQ
type Message struct {
	Method    string
	Name      string
	Email     string
	Activated string
	Sam       string
	Password  string
}

// Plugin structure
type Ldap struct{}

// Strucutre used in return messages sent to RabbitMQ
type ReturnMsg struct {
	Method string
	Err    string
	Plugin string
	Email  string
}

func env(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func (h *handler) listUsers(req nano.Request) (*nano.Response, error) {

	res, err := h.ldapCon.GetUsers()
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil

	}

	return nano.JSONResponse(200, res), nil
}

// Checks if there is at least one sam account available, to use it to create a new user instead of generating a new sam account

func (h *handler) updatePassword(req nano.Request) (*nano.Response, error) {
	var params struct {
		UserEmail string
		Password  string
	}
	err := json.Unmarshal(req.Body, &params)
	if err != nil {
		module.Log.Error("Unable to unmarshall params: " + err.Error())
		return nil, err
	}

	if len(req.Params["user_id"]) < 1 {
		return nano.JSONResponse(400, hash{
			"error": "user id is missing",
		}), nil
	}

	params.UserEmail = req.Params["user_id"]

	err = h.ldapCon.ChangePassword(params.UserEmail, params.Password)
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil

	}
	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func (h *handler) createUser(req nano.Request) (*nano.Response, error) {
	var params struct {
		UserEmail string
		Password  string
	}

	err := json.Unmarshal(req.Body, &params)

	if err != nil {
		module.Log.Error("Unable to unmarshall params: " + err.Error())
		return nano.JSONResponse(400, hash{
			"error": err.Error(),
		}), err
	}

	if params.UserEmail == "" || params.Password == "" {
		module.Log.Error("Email or password missing")
		return nano.JSONResponse(400, hash{
			"error": "Email or password missing",
		}), nil
	}

	sam, err := h.ldapCon.AddUser(params.UserEmail, params.Password)
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}
	return nano.JSONResponse(200, hash{
		"sam": sam,
	}), nil
}

func (h *handler) forcedisableAccount(req nano.Request) (*nano.Response, error) {
	userId := req.Params["user_id"]

	if len(userId) < 1 {
		module.Log.Error("User ID missing")
		return nano.JSONResponse(400, hash{
			"error": "User id is missing",
		}), nil
	}

	err := h.ldapCon.DisableUser(userId)
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func (h *handler) deleteUser(req nano.Request) (*nano.Response, error) {
	userId := req.Params["user_id"]

	if len(userId) < 1 {
		module.Log.Error("User ID missing")
		return nano.JSONResponse(400, hash{
			"error": "User id is missing",
		}), nil
	}

	err := h.ldapCon.DeleteAccount(userId)
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func main() {
	conf.Username = env("LDAP_USERNAME", "CN=Administrator,CN=Users,DC=intra,DC=localdomain,DC=com")
	conf.Password = env("LDAP_PASSWORD", "Nanocloud123+")
	conf.Ou = env("LDAP_OU", "OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com")

	ldapServer, err := url.Parse(env("LDAP_SERVER_URI", "ldaps://Administrator:Nanocloud123+@172.17.0.1:6003"))
	if err != nil {
		panic(err)
	}

	conf.LDAPServer = *ldapServer

	module = nano.RegisterModule("ldap")

	h := handler{
		ldapCon: ldapPkg.New(conf.Username, conf.Password, conf.ServerURL, conf.Ou, conf.LDAPServer),
	}

	module.Post("/ldap/users", h.createUser)
	module.Get("/ldap/users", h.listUsers)
	module.Put("/ldap/users/:user_id", h.updatePassword)
	module.Post("/ldap/users/:user_id/disable", h.forcedisableAccount)
	module.Delete("/ldap/users/:user_id", h.deleteUser)

	module.Listen()
}
