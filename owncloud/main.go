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
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"strings"

	"github.com/natefinch/pie"

	//todo vendor this dependency
	// nan "nanocloud.com/plugins/owncloud/libnan"
)

// Create an object to be exported

var (
	name = "owncloud"
	srv  pie.Server
)

type CreateUserParams struct {
	Username, Password string
}

type api struct{}

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
}

func CreateUser(args PlugRequest, reply *PlugRequest) error {
	var params CreateUserParams
	err := json.Unmarshal([]byte(args.Body), &params)
	if err != nil {
		log.Println(err)
	}
	_, err = Create(params.Username, params.Password)
	if err != nil {
		log.Println(err)
	}
	return err
}

func ChangePassword(args PlugRequest, reply *PlugRequest) {
	var params CreateUserParams
	err := json.Unmarshal([]byte(args.Body), &params)
	if err != nil {
		log.Println(err)
	}
	_, err = Edit(params.Username, "password", params.Password)
}

type del struct {
	Username string
}

func DeleteUser(args PlugRequest, reply *PlugRequest) {

	var User del
	err := json.Unmarshal([]byte(args.Body), &User)
	if err != nil {
		log.Println(err)
	}
	_, err = Delete(User.Username)
	if err != nil {
		log.Println("deletion error: ", err)
	}
}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	initConf()
	Configure()

	if strings.Index(args.Url, "/owncloud/add") == 0 {
		CreateUser(args, reply)
	}
	if strings.Index(args.Url, "/owncloud/delete") == 0 {
		DeleteUser(args, reply)
	}
	if strings.Index(args.Url, "/owncloud/changepassword") == 0 {
		ChangePassword(args, reply)
	}

	return nil
}

func (api) Plug(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Unplug(args interface{}, reply *bool) error {
	defer os.Exit(0)
	*reply = true
	return nil
}

func main() {
	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)

}
