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
	"github.com/Nanocloud/nano"
	"github.com/Nanocloud/oauth"
)

type oauthConnector struct{}

type Client struct {
	Id   int
	Name string
	Key  string
}

type AccessToken struct {
	Token string
	Type  string
}

func (c oauthConnector) AuthenticateUser(username, password string) (interface{}, error) {
	m := make(map[string]string)
	m["username"] = username
	m["password"] = password

	res, err := module.JSONRequest("POST", "/users/login", m, nil)
	if err != nil {
		module.Log.Error(err)
		return nil, err
	}

	module.Log.Debug(res.StatusCode)
	if res.StatusCode != 200 {
		return nil, nil
	}

	user := nano.User{}

	module.Log.Debug(string(res.Body))
	err = json.Unmarshal(res.Body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c oauthConnector) GetUserFromAccessToken(accessToken string) (interface{}, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`SELECT user_id
		FROM oauth_access_tokens
		WHERE token = $1::varchar`,
		accessToken,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if !rows.Next() {
		return nil, nil
	}

	var userId string
	err = rows.Scan(&userId)
	if err != nil {
		return nil, err
	}

	res, err := module.Request("GET", "/users/"+userId, "", nil, nil)
	if err != nil {
		module.Log.Error(err)
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, nil
	}

	user := nano.User{}

	err = json.Unmarshal(res.Body, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (c oauthConnector) GetClient(key string, secret string) (interface{}, error) {
	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`SELECT id, name,
		key
		FROM oauth_clients
		WHERE key = $1::varchar
		AND secret = $2::varchar`,
		key, secret,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		client := Client{}

		err = rows.Scan(&client.Id, &client.Name, &client.Key)
		return &client, nil
	}
	return nil, nil
}

func (at AccessToken) ToJSON() ([]byte, error) {
	m := make(map[string]string)
	m["access_token"] = at.Token
	m["type"] = at.Type

	return json.Marshal(&m)
}

func (c oauthConnector) GetAccessToken(rawUser, rawClient interface{}) (oauth.JSONAble, error) {
	user := rawUser.(*nano.User)
	client := rawClient.(*Client)

	db, err := GetDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(
		`SELECT token
		FROM oauth_access_tokens
		WHERE oauth_client_id = $1::integer
		AND user_id = $2::varchar`,
		client.Id, user.Id,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	if rows.Next() {
		accessToken := AccessToken{}
		accessToken.Type = "Bearer"
		rows.Scan(&accessToken.Token)
		return accessToken, nil
	}

	token := randomString(25)

	rows, err = db.Query(
		`INSERT INTO oauth_access_tokens
		(token, oauth_client_id, user_id)
		VALUES ($1::varchar, $2::integer, $3::varchar)`,
		token, client.Id, user.Id,
	)

	if err != nil {
		return nil, err
	}

	rows.Close()

	accessToken := AccessToken{token, "Bearer"}
	return accessToken, nil
}

func init() {
	oauth.SetConnector(oauthConnector{})
}
