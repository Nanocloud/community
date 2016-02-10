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
	"errors"
	"github.com/Nanocloud/community/nanocloud/connectors/db"
	log "github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	UserNotFound       = errors.New("user not found")
	InvalidCredentials = errors.New("invalid credentials")
	UserDisabled       = errors.New("user disabled")
	UserDuplicated     = errors.New("user duplicated")
	UserNotCreated     = errors.New("user not created")
)

func GetUserFromEmailPassword(email, password string) (*User, error) {
	if len(email) < 1 || len(password) < 1 {
		return nil, UserNotFound
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
		return nil, err
	}

	if !rows.Next() {
		return nil, UserNotFound
	}

	var user User
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

func FindUsers() (*[]User, error) {
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

	var users []User

	defer rows.Close()
	for rows.Next() {
		user := User{}

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

func UserExists(id string) (bool, error) {
	rows, err := db.Query(
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

func DisableUser(id string) error {
	rows, err := db.Query(
		`UPDATE users
		SET activated = false
		WHERE id = $1::varchar`,
		id)
	if err != nil {
		rows.Close()
	}
	return err
}

func CreateUser(
	activated bool,
	email string,
	firstName string,
	lastName string,
	password string,
	isAdmin bool,
) (*User, error) {
	id := uuid.NewV4().String()

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
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

	rows, err = db.Query(
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

	var user User
	rows.Scan(
		&user.Id, &user.Activated,
		&user.Email, &user.FirstName,
		&user.LastName, &user.IsAdmin,
		&user.Sam, &user.WindowsPassword,
	)

	rows.Close()

	return &user, err
}

func UpdateUserAd(id, sam, password string) error {
	res, err := db.Exec(
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

func DeleteUser(id string) error {
	res, err := db.Exec("DELETE FROM users WHERE id = $1::varchar", id)
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

func UpdateUserPassword(id string, password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	res, err := db.Exec(
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

func GetUser(id string) (*User, error) {
	rows, err := db.Query(
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
		var user User

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
