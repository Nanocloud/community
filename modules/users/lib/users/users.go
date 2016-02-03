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

package users

import (
	"database/sql"
	"errors"
	"github.com/Nanocloud/nano"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	UserNotFound       = errors.New("user not found")
	InvalidCredentials = errors.New("invalid credentials")
	UserDisabled       = errors.New("user disabled")
	UserDuplicated     = errors.New("user duplicated")
	UserNotCreated     = errors.New("user not created")
)

type Users struct {
	db *sql.DB
}

func dbConnect(databaseURI string) (*sql.DB, error) {
	var err error
	var db *sql.DB

	for try := 0; try < 10; try++ {
		db, err = sql.Open("postgres", databaseURI)
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err == nil {
			log.Info("Connected to Postgres")
			return db, nil
		}

		log.Info("Unable to connect to Postgres. Will retry in 5 sec")
		time.Sleep(time.Second * 5)
	}
	return nil, err
}

func (u *Users) GetUserFromEmailPassword(email, password string) (*nano.User, error) {
	if len(email) < 1 || len(password) < 1 {
		return nil, UserNotFound
	}

	rows, err := u.db.Query(
		`SELECT id, activated,
		email, password,
		first_name, last_name,
		is_admin
		FROM users
		WHERE email = $1::varchar`,
		email,
	)
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, UserNotFound
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
		return nil, InvalidCredentials
	}

	if !user.Activated {
		return nil, UserDisabled
	}

	return &user, nil
}

func (u *Users) FindUsers() (*[]nano.User, error) {
	rows, err := u.db.Query(
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

func (u *Users) UserExists(id string) (bool, error) {
	rows, err := u.db.Query(
		`SELECT id
		FROM users
		WHERE id = $1::varchar`,
		id)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	if rows.Next() {
		return true, nil
	}
	return false, nil
}

func (u *Users) DisableUser(id string) error {
	rows, err := u.db.Query(
		`UPDATE users
		SET activated = false
		WHERE id = $1::varchar`,
		id)
	if err != nil {
		rows.Close()
	}
	return err
}

func (u *Users) CreateUser(
	activated bool,
	email string,
	firstName string,
	lastName string,
	password string,
	isAdmin bool,
) (*nano.User, error) {
	id := uuid.NewV4().String()

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rows, err := u.db.Query(
		`INSERT INTO users
		(id, email, activated,
		first_name, last_name,
		password, is_admin)
		VALUES(
			$1::varchar, $2::varchar, $3::bool,
			$4::varchar, $5::varchar,
			$6::varchar, $7::bool)
		`, id, email, activated,
		firstName, lastName,
		pass, isAdmin)

	if err != nil {
		switch err.Error() {
		case "pq: duplicate key value violates unique constraint \"users_pkey\"":
			log.Error("user id exists already")
			return nil, UserDuplicated
		case "pq: duplicate key value violates unique constraint \"users_email_key\"":
			log.Error("user email exists already")
			return nil, UserDuplicated
		}
		return nil, err
	}

	rows.Close()

	rows, err = u.db.Query(
		`SELECT id, activated,
		email,
		first_name, last_name,
		is_admin, sam, windows_password
		FROM users
		WHERE id = $1::varchar`,
		id)

	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, UserNotCreated
	}

	var user nano.User
	rows.Scan(
		&user.Id, &user.Activated,
		&user.Email, &user.FirstName,
		&user.LastName, &user.IsAdmin,
		&user.Sam, &user.WindowsPassword,
	)

	rows.Close()

	return &user, err
}

func (u *Users) UpdateUserAd(id, sam, password string) error {
	res, err := u.db.Exec(
		`UPDATE users
		SET sam = $1::varchar,
		windows_password = $2::varchar
		WHERE id = $3::varchar`,
		sam, password, id)
	if err != nil {
		return err
	}
	updated, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return UserNotFound
	}
	return nil
}

func (u *Users) DeleteUser(id string) error {
	res, err := u.db.Exec("DELETE FROM users WHERE id = $1::varchar", id)
	if err != nil {
		return err
	}
	deleted, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if deleted == 0 {
		return UserNotFound
	}
	return nil
}

func (u *Users) UpdateUserPassword(id string, password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := u.db.Exec(
		`UPDATE users
		SET password = $1::varchar
		WHERE id = $2::varchar`,
		pass, id)
	if err != nil {
		return err
	}

	updated, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if updated == 0 {
		return UserNotFound
	}
	return nil
}

func (u *Users) GetUser(id string) (*nano.User, error) {
	rows, err := u.db.Query(
		`SELECT id,
		first_name, last_name,
		email, is_admin, activated,
		sam, windows_password
		FROM users
		WHERE id = $1::varchar`,
		id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		var user nano.User

		err = rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.IsAdmin,
			&user.Activated,
			&user.Sam,
			&user.WindowsPassword,
		)
		if err != nil {
			return nil, err
		}

		return &user, nil
	}
	return nil, nil
}

func (u *Users) init() error {
	rows, err := u.db.Query(
		`SELECT table_name
			FROM information_schema.tables
			WHERE table_name = 'users'`)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		return nil
	}

	rows, err = u.db.Query(
		`CREATE TABLE users (
				id               varchar(36) PRIMARY KEY,
				first_name       varchar(36) NOT NULL DEFAULT '',
				last_name        varchar(36) NOT NULL DEFAULT '',
				email            varchar(36) NOT NULL DEFAULT '' UNIQUE,
				password         varchar(60) NOT NULL DEFAULT '',
				is_admin         boolean,
				activated        boolean,
				sam              varchar(35) NOT NULL DEFAULT '',
				windows_password varchar(36) NOT NULL DEFAULT ''
			);`)
	if err != nil {
		return err
	}

	rows.Close()

	_, err = u.CreateUser(
		true,
		"admin@nanocloud.com",
		"John",
		"Doe",
		"admin",
		true,
	)

	if err != nil {
		return err
	}
	return nil
}

func New(dbURI string) *Users {
	db, err := dbConnect(dbURI)
	if err != nil {
		panic(err)
	}

	u := Users{
		db: db,
	}
	err = u.init()
	if err != nil {
		panic(err)
	}

	return &u
}
