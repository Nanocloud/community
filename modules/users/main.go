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
	"math/rand"
	"os"
	"time"

	"github.com/Nanocloud/nano"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type hash map[string]interface{}

var module nano.Module

var db *sql.DB

func dbConnect() error {
	databaseURI := os.Getenv("DATABASE_URI")
	if len(databaseURI) == 0 {
		databaseURI = "postgres://localhost/nanocloud?sslmode=disable"
	}

	var err error

	for try := 0; try < 10; try++ {
		db, err = sql.Open("postgres", databaseURI)
		if err != nil {
			return err
		}

		err = db.Ping()
		if err == nil {
			module.Log.Info("Connected to Postgres")
			return nil
		}

		module.Log.Info("Unable to connect to Postgres. Will retry in 5 sec")
		time.Sleep(time.Second * 5)
	}

	return err
}

func findUsers() (*[]nano.User, error) {
	rows, err := db.Query(
		`SELECT id,
		first_name, last_name,
		email, is_admin, activated,
		sam, windows_password
		FROM users`,
	)
	if err != nil {
		return nil, err
	}

	var users []nano.User

	defer rows.Close()
	for rows.Next() {
		user := nano.User{}

		rows.Scan(
			&user.Id,
			&user.FirstName, &user.LastName,
			&user.Email,
			&user.IsAdmin,
			&user.Activated,
			&user.Sam,
			&user.WindowsPassword,
		)
		users = append(users, user)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &users, nil
}

func getUser(req nano.Request) (*nano.Response, error) {
	userId := req.Params["id"]
	if userId == "" {
		return nano.JSONResponse(400, hash{
			"error": "User id needed to retrieve account informations",
		}), nil
	}

	rows, err := db.Query(
		`SELECT id,
		first_name, last_name,
		email, is_admin, activated,
		sam, windows_password
		FROM users
		WHERE id = $1::varchar`,
		userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		var user nano.User
		rows.Scan(
			&user.Id,
			&user.FirstName, &user.LastName,
			&user.Email,
			&user.IsAdmin,
			&user.Activated,
			&user.Sam,
			&user.WindowsPassword,
		)
		return nano.JSONResponse(200, user), nil
	}

	return nano.JSONResponse(404, hash{
		"error": "User Not Found",
	}), nil
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

func CreateADUser(id string) (string, string, error) {
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

func CreateUser(
	activated bool,
	email string,
	firstName string,
	lastName string,
	password string,
	isAdmin bool,
) (createdUser *nano.User, err error) {
	id := uuid.NewV4().String()
	sam, winpass, err := CreateADUser(id)
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	rows, err := db.Query(
		`INSERT INTO users
		(id, email, activated,
		first_name, last_name,
		password, is_admin,
		sam, windows_password)
		VALUES(
			$1::varchar, $2::varchar, $3::bool,
			$4::varchar, $5::varchar,
			$6::varchar, $7::bool,
			$8::varchar, $9::varchar)
		`, id, email, activated,
		firstName, lastName,
		pass, isAdmin, sam, winpass)

	if err != nil {
		switch err.Error() {
		case "pq: duplicate key value violates unique constraint \"users_pkey\"":
			err = errors.New("user id exists already")
		case "pq: duplicate key value violates unique constraint \"users_email_key\"":
			err = errors.New("user email exists already")
		}
		return
	}

	rows.Close()

	rows, err = db.Query(
		`SELECT id, activated,
		email,
		first_name, last_name,
		is_admin, sam, windows_password
		FROM users
		WHERE id = $1::varchar`,
		id)

	if err != nil {
		return
	}

	if !rows.Next() {
		err = errors.New("user not created")
		return
	}

	var user nano.User
	rows.Scan(
		&user.Id, &user.Activated,
		&user.Email, &user.FirstName,
		&user.LastName, &user.IsAdmin,
		&user.Sam, &user.WindowsPassword,
	)

	rows.Close()

	createdUser = &user
	return
}

func updateUserPassword(req nano.Request) (*nano.Response, error) {
	userId := req.Params["id"]
	if userId == "" {
		return nano.JSONResponse(400, hash{
			"error": "Email needed to modify account",
		}), nil
	}

	var user struct {
		Password string
	}

	err := json.Unmarshal(req.Body, &user)
	if err != nil {
		return nil, err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`UPDATE users
		SET password = $1::varchar
		WHERE id = $2::varchar`,
		pass, userId)
	if err != nil {
		return nil, err
	}

	rows.Close()

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

	rows, err := db.Query(
		`UPDATE users
		SET activated = false
		WHERE id = $1::varchar`,
		userId)

	if err != nil {
		return nil, err
	}
	rows.Close()

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func deleteUser(req nano.Request) (*nano.Response, error) {
	userId := req.Params["id"]
	if len(userId) == 0 {
		return nano.JSONResponse(400, hash{
			"error": "User id needed for deletion",
		}), nil
	}

	rows, err := db.Query("DELETE FROM users WHERE id = $1::varchar", userId)
	if err != nil {
		module.Log.Error(err)
		return nil, err
	}
	rows.Close()
	// SendMsg(Message{Method: "Delete", Email: mail})

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func getUsers(req nano.Request) (*nano.Response, error) {
	users, err := findUsers()
	if err != nil {
		return nil, err
	}
	return nano.JSONResponse(200, users), nil
}

func getUserFromEmailPassword(email, password string) (*nano.User, string, error) {
	if len(email) < 1 || len(password) < 1 {
		return nil, "user not found", nil
	}

	rows, err := db.Query(
		`SELECT id, activated,
		email, password,
		first_name, last_name,
		is_admin
		FROM users
		WHERE email = $1::varchar`,
		email,
	)
	if err != nil {
		return nil, "", err
	}

	if !rows.Next() {
		return nil, "user not found", nil
	}

	var user nano.User
	var passwordHash string
	rows.Scan(
		&user.Id, &user.Activated,
		&user.Email, &passwordHash,
		&user.FirstName, &user.LastName,
		&user.IsAdmin,
	)
	rows.Close()

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))

	if err != nil {
		return nil, "wrong password", nil
	}

	if !user.Activated {
		return nil, "user is not activated", nil
	}

	return &user, "", nil
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

	user, message, err := getUserFromEmailPassword(body.Username, body.Password)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return nano.JSONResponse(200, hash{
			"success": true,
			"user":    user,
		}), nil
	}

	return nano.JSONResponse(400, hash{
		"success": false,
		"error":   message,
	}), nil
}

func setupDb() error {
	rows, err := db.Query(
		`SELECT table_name
			FROM information_schema.tables
			WHERE table_name = 'users'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		module.Log.Info("Users table already set up")
		return nil
	}

	rows, err = db.Query(
		`CREATE TABLE users (
				id               varchar(36) PRIMARY KEY,
				first_name       varchar(36),
				last_name        varchar(36),
				email            varchar(36) UNIQUE,
				password         varchar(60),
				is_admin         boolean,
				activated        boolean,
				sam        	 varchar(35),
				windows_password varchar(36)
			);`)
	if err != nil {
		module.Log.Errorf("Unable to create users table: %s", err)
		return err
	}

	rows.Close()

	_, err = CreateUser(
		true,
		"admin@nanocloud.com",
		"John",
		"Doe",
		"admin",
		true,
	)

	if err != nil {
		module.Log.Errorf("Unable to create the default user: %s", err)
		return err
	}
	return nil
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

	_, err = CreateUser(
		true,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		false,
	)
	if err != nil {
		return nil, err
	}

	return nano.JSONResponse(200, hash{
		"success": true,
	}), nil
}

func main() {
	module = nano.RegisterModule("users")

	err := dbConnect()
	if err != nil {
		module.Log.Fatalf("Cannot connect to Postgres Database: %s", err)
	}

	err = setupDb()
	if err != nil {
		module.Log.Fatalf("Unable to setup users table: %s", err)
		return
	}

	module.Post("/users/login", userLogin)

	module.Post("/users/:id/disable", disableUser)
	module.Get("/users", getUsers)

	// Create a User
	module.Post("/users", postUsers)

	module.Delete("/users/:id", deleteUser)
	module.Put("/users/:id", updateUserPassword)
	module.Get("/users/:id", getUser)

	module.Listen()
}
