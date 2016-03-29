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

package oauth

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/oauth2"
	"github.com/Nanocloud/community/nanocloud/utils"
	"github.com/satori/go.uuid"
)

type oauthConnector struct{}

type Client struct {
	Id   int
	Name string
	Key  string
}

type AccessToken struct {
	Token     string        `json:"access_token"`
	Type      string        `json:"token_type"`
	ExpiresIn time.Duration `json:"expires_in"`
}

func (c oauthConnector) AuthenticateUser(username, password string) (interface{}, error) {
	return users.GetUserFromEmailPassword(username, password)
}

func (c oauthConnector) GetUserFromAccessToken(accessToken string) (interface{}, error) {
	rows, err := db.Query(
		`SELECT user_id
		FROM oauth_access_tokens
		WHERE token = $1::varchar
		AND expires_at > NOW()
		`,
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

	return users.GetUser(userId)
}

func (c oauthConnector) GetClient(key string, secret string) (interface{}, error) {
	rows, err := db.Query(
		`SELECT id, name,
		key
		FROM oauth_clients
		WHERE key = $1::varchar`,
		key,
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

func removeExpiredTokens() {
	db.Exec(
		`DELETE FROM oauth_access_tokens
		WHERE expires_at < NOW()`,
	)
}

func (c oauthConnector) GetAccessToken(rawUser, rawClient interface{}, req *http.Request) (interface{}, error) {
	removeExpiredTokens()

	user := rawUser.(*users.User)
	client := rawClient.(*Client)

	ua := req.UserAgent()

	// Get IP client address
	var ip string
	if os.Getenv("TRUST_PROXY") == "true" {
		xForwardedFor := req.Header["X-Forwarded-For"]
		if len(xForwardedFor) > 0 {
			ip = xForwardedFor[0]
		}
	}

	if len(ip) == 0 {
		addr := req.RemoteAddr
		i := strings.LastIndex(addr, ":")
		ip = addr[0:i]
	}

	id := uuid.NewV4().String()
	token := utils.RandomString(25)

	rows, err := db.Query(
		`INSERT INTO oauth_access_tokens
		(id, token, oauth_client_id, user_id,
		 created_at, user_agent, ip, expires_at)
		VALUES
		($1::varchar, $2::varchar, $3::integer, $4::varchar,
		 NOW(), $5::varchar, $6::varchar, NOW() + interval '1 day')`,
		id, token, client.Id, user.Id,
		ua, ip,
	)

	if err != nil {
		return nil, err
	}

	rows.Close()

	accessToken := AccessToken{
		Token:     token,
		Type:      "Bearer",
		ExpiresIn: 60 * 60 * 24,
	}
	return accessToken, nil
}

func init() {
	oauth2.SetConnector(oauthConnector{})
}
