package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/nanocloud/oauth"
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
	userModule := conf.RunDir + "users"

	args := struct {
		Username string
		Password string
	}{
		username,
		password,
	}

	res := struct {
		Success      bool
		ErrorMessage string
		User         UserInfo
	}{}

	err := plugins[userModule].client.Call("users.AuthenticateUser", args, &res)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if !res.Success {
		return nil, nil
	}

	return &res.User, nil
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

	userModule := conf.RunDir + "users"

	args := struct {
		UserId string
	}{
		userId,
	}

	res := struct {
		Success      bool
		ErrorMessage string
		User         UserInfo
	}{}

	err = plugins[userModule].client.Call("users.GetUser", args, &res)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if !res.Success {
		return nil, nil
	}

	return &res.User, nil
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
	log.Debug(key)
	log.Debug(secret)

	defer rows.Close()
	if rows.Next() {
		client := Client{}

		err = rows.Scan(&client.Id, &client.Name, &client.Key)
		if err != nil {
			log.Debug(err)
		}
		log.Debugf("Client id = %d\n", client.Id)
		log.Debug(client.Name)
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
	user := rawUser.(*UserInfo)
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

	log.Debugf("Client id = %d\n", client.Id)
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
