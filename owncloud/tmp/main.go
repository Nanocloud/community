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
	"fmt"
	"log"

	//todo vendor this dependency
	// nan "nanocloud.com/plugins/owncloud/libnan"
)

// Create an object to be exported

type CreateUserParams struct {
	Username, Password string
}

func Configure(jsonConfig string, pOutMsg *string) error {

	var params map[string]string
	if err := json.Unmarshal([]byte(jsonConfig), &params); err != nil {
		*pOutMsg = fmt.Sprintf("Configure() failed to unmarshal ownCloud plugin configuration: %s", err)
		return err
	}
	Config(params["protocol"], params["url"], params["login"], params["password"])

	return nil
}

func CreateUser(params CreateUserParams, reply *bool) error {
	_, err := Create(params.Username, params.Password)
	*reply = err == nil
	return err
}

func ChangePassword(params CreateUserParams, reply *bool) error {
	_, err := Edit(params.Username, "password", params.Password)
	*reply = err == nil
	return err
}

func DeleteUser(username string, reply *bool) error {
	_, err := Delete(username)
	*reply = err == nil
	return err
}

func main() {
	out := ""
	err := Configure(`{ "protocol" : "https", "url" : "192.168.1.39/drive", "login" : "drive_admin", "password" : "BJboVHDiawECoDt" }`, &out)
	if err != nil {
		log.Println(err)
	}
	var params CreateUserParams
	params.Username = "joe2"
	params.Password = "Nanocloud123+"
	/*b := false
	err = CreateUser(params, &b)
	if err != nil {
		log.Println("Create user error: ", err)
	}*/
	res, err := List()
	if err != nil {
		log.Println("list error :", err)
	}
	log.Println(res)
}
