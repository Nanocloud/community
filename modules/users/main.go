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
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"time"

	usersPkg "github.com/Nanocloud/community/modules/users/lib/users"
	"github.com/Nanocloud/nano"
)

type hash map[string]interface{}

var module nano.Module
var users *usersPkg.Users

func AdminOnly(req nano.Request) (*nano.Response, error) {
	if req.User != nil && !req.User.IsAdmin {
		return nano.JSONResponse(403, hash{
			"error": "forbidden",
		}), nil
	}
	return nil, nil
}

func getUser(req nano.Request) (*nano.Response, error) {
	userId := req.Params["id"]
	if userId == "" {
		return nano.JSONResponse(400, hash{
			"error": "User id needed to retrieve account informations",
		}), nil
	}

	user, err := users.GetUser(userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nano.JSONResponse(404, hash{
			"error": "User Not Found",
		}), nil
	}

	return nano.JSONResponse(200, user), nil
}

// randomString
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randomString(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func createADUser(id string) (string, string, error) {
	password := randomString(8) + "s4D+"
	res, err := module.JSONRequest("POST", "/ldap/users", hash{
		"userEmail": id,
		"password":  password,
	}, nil)
	if err != nil {
		return "", "", err
	}
	var r struct {
		Sam string
	}
	err = json.Unmarshal(res.Body, &r)
	if err != nil {
		return "", "", err
	}
	return r.Sam, password, nil
}

func updateUserPassword(req nano.Request) (*nano.Response, error) {
	userId := req.Params["id"]
	if userId == "" {
		return nano.JSONResponse(400, hash{
			"error": "User id needed to modify account",
		}), nil
	}

	var user struct {
		Password string
	}

	err := json.Unmarshal(req.Body, &user)
	if err != nil {
		module.Log.Errorf("Unable to parse body request: %s", err.Error())
		return nil, err
	}

	exists, err := users.UserExists(userId)
	if err != nil {
		module.Log.Errorf("Unable to check user existance: %s", err.Error())
		return nil, err
	}

	if !exists {
		return nano.JSONResponse(404, hash{
			"error": "User not found",
		}), nil
	}

	err = users.UpdateUserPassword(userId, user.Password)
	if err != nil {
		module.Log.Errorf("Unable to update user password: %s", err.Error())
		return nil, err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func disableUser(req nano.Request) (*nano.Response, error) {
	userId := req.Params["id"]
	if userId == "" {
		return nano.JSONResponse(404, hash{
			"error": "User id needed for desactivation",
		}), nil
	}

	exists, err := users.UserExists(userId)
	if err != nil {
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	if !exists {
		return nano.JSONResponse(404, hash{
			"error": "User not found",
		}), nil
	}

	err = users.DisableUser(userId)
	if err != nil {
		module.Log.Errorf("Unable to disable user: %s", err.Error())
		return nil, err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func deleteADUser(id string) error {
	_, err := module.JSONRequest("DELETE", "/ldap/users/"+id, hash{}, nil)
	return err
}

func deleteUser(req nano.Request) (*nano.Response, error) {
	userId := req.Params["id"]
	if len(userId) == 0 {
		return nano.JSONResponse(400, hash{
			"error": "User id needed for deletion",
		}), nil
	}

	user, err := users.GetUser(userId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nano.JSONResponse(404, hash{
			"error": "User not found",
		}), nil
	}

	if user.IsAdmin {
		return nano.JSONResponse(403, hash{
			"error": "Admins cannot be deleted",
		}), nil
	}

	err = deleteADUser(user.Id)
	if err != nil {
		module.Log.Errorf("Unable to delete user in ad: %s", err.Error())
		return nil, err
	}

	err = users.DeleteUser(user.Id)
	if err != nil {
		module.Log.Errorf("Unable to delete user: ", err.Error())
		return nil, err
	}

	// SendMsg(Message{Method: "Delete", Email: mail})

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func getUsers(req nano.Request) (*nano.Response, error) {
	users, err := users.FindUsers()
	if err != nil {
		module.Log.Errorf("unable to get user lists: %s", err.Error())
		return nil, err
	}
	return nano.JSONResponse(200, users), nil
}

func userLogin(req nano.Request) (*nano.Response, error) {
	var body struct {
		Username string
		Password string
	}

	err := json.Unmarshal(req.Body, &body)
	if err != nil {
		return nil, err
	}

	user, err := users.GetUserFromEmailPassword(body.Username, body.Password)
	switch err {
	case usersPkg.InvalidCredentials:
	case usersPkg.UserDisabled:
		return nano.JSONResponse(400, hash{
			"success": false,
			"error":   err.Error(),
		}), nil

	case nil:
		if user == nil {
			module.Log.Error("unable to log the user in")
			return nil, errors.New("unable to log the user in")
		}

		return nano.JSONResponse(200, hash{
			"success": true,
			"user":    user,
		}), nil
	}
	return nil, err
}

func postUsers(req nano.Request) (*nano.Response, error) {
	var user struct {
		Email     string
		FirstName string
		LastName  string
		Password  string
	}

	err := json.Unmarshal([]byte(req.Body), &user)
	if err != nil {
		module.Log.Error(err)
		return nil, err
	}

	newUser, err := users.CreateUser(
		true,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		false,
	)

	switch err {
	case usersPkg.UserDuplicated:
		return nano.JSONResponse(409, hash{
			"error": err.Error(),
		}), nil
	case usersPkg.UserNotCreated:
		return nano.JSONResponse(500, hash{
			"error": err.Error(),
		}), nil
	}

	sam, winpass, err := createADUser(newUser.Id)
	err = users.UpdateUserAd(newUser.Id, sam, winpass)

	if err != nil {
		return nil, err
	}

	return nano.JSONResponse(201, hash{
		"Id": newUser.Id,
	}), nil
}

func main() {
	module = nano.RegisterModule("users")

	databaseURI := os.Getenv("DATABASE_URI")
	if len(databaseURI) == 0 {
		databaseURI = "postgres://localhost/nanocloud?sslmode=disable"
	}

	users = usersPkg.New(databaseURI)

	module.Post("/users/login", AdminOnly, userLogin)

	module.Post("/users/:id/disable", AdminOnly, disableUser)
	module.Get("/users", AdminOnly, getUsers)

	// Create a User
	module.Post("/users", AdminOnly, postUsers)

	module.Delete("/users/:id", AdminOnly, deleteUser)
	module.Put("/users/:id", AdminOnly, updateUserPassword)
	module.Get("/users/:id", AdminOnly, getUser)

	module.Listen()
}
