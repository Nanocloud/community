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
	"fmt"
	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
	"github.com/natefinch/pie"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"regexp"
)

var (
	name = "users"
	srv  pie.Server
)

type hash map[string]interface{}

type api struct{}

var db *sql.DB

type UserInfo struct {
	Id              string
	Activated       bool
	Email           string
	FirstName       string
	LastName        string
	Password        string
	IsAdmin         bool
	Sam             string
	WindowsPassword string
}

type Message struct {
	Method    string
	Name      string
	Email     string
	Activated string
	Sam       string
	Password  string
}

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
	Method   string
	HeadVals map[string]string
	Status   int
}

type ReturnMsg struct {
	Method string
	Err    string
	Plugin string
	Email  string
}

func GetUsers() (*[]UserInfo, error) {
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

	var users []UserInfo

	defer rows.Close()
	for rows.Next() {
		user := UserInfo{}

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

func GetUser(args PlugRequest, reply *PlugRequest, userId string) (err error) {
	reply.Status = 400
	if userId == "" {
		err = errors.New("User id needed to retrieve account informations")

		log.Error(err)
		return
	}

	reply.Status = 500
	rows, err := db.Query(
		`SELECT id,
		first_name, last_name,
		email, is_admin, activated,
		sam, windows_password
		FROM users
		WHERE id = $1::varchar`,
		userId)
	if err != nil {
		return
	}

	defer rows.Close()
	if rows.Next() {
		reply.Status = 200

		reply.HeadVals = make(map[string]string, 1)
		reply.HeadVals["Content-Type"] = "application/json; charset=UTF-8"

		var user UserInfo
		rows.Scan(
			&user.Id,
			&user.FirstName, &user.LastName,
			&user.Email,
			&user.IsAdmin,
			&user.Activated,
			&user.Sam,
			&user.WindowsPassword,
		)

		var res []byte
		res, err = json.Marshal(user)
		if err != nil {
			reply.Status = 500
			return
		}

		reply.Status = 200
		reply.Body = string(res)
	} else {
		reply.Status = 404
		err = errors.New("User Not Found")
	}

	return
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Error("%s: %s", msg, err)
	}
}

func SendMsg(msg Message) {
	conn, err := amqp.Dial(conf.QueueUri)
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	err = ch.ExchangeDeclare(
		"users_topic", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	str, err := json.Marshal(msg)
	if err != nil {
		log.Error(err)
	}
	err = ch.Publish(
		"users_topic", // exchange
		"users.req",   // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(str),
		})
	failOnError(err, "Failed to publish a message")

	log.Info(" [x] Sent order to plugin")
	defer ch.Close()
	defer conn.Close()

}

func CreateADUser(id string) (string, string, error) {
	password := randomString(8) + "s4D+"
	args := make(map[string]string, 2)
	args["id"] = id
	args["password"] = password
	log.Error("CALLING RPCREQUEST")
	res, err := rpcRequest("rmq_ldap", "create_user", args)
	log.Error("CALLED RPCREQUEST")
	return res["sam"].(string), password, err
}

func CreateUser(
	activated bool,
	email string,
	firstName string,
	lastName string,
	password string,
	isAdmin bool,
) (createdUser *UserInfo, err error) {
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

	var user UserInfo
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

func Add(args PlugRequest, reply *PlugRequest, mail string) (err error) {
	var user UserInfo
	err = json.Unmarshal([]byte(args.Body), &user)
	if err != nil {
		log.Error(err)
		return
	}

	_, err = CreateUser(
		true,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Password,
		false,
	)
	if err == nil {
		reply.Status = 202
	} else {
		reply.Status = 400
	}
	return
}

func UpdatePassword(args PlugRequest, reply *PlugRequest, userId string) (err error) {
	reply.Status = 400
	if userId == "" {
		err = errors.New("Email needed to modify account")
		return
	}
	reply.Status = 500

	var user UserInfo
	err = json.Unmarshal([]byte(args.Body), &user)
	if err != nil {
		log.Error(err)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	rows, err := db.Query(
		`UPDATE users
		SET password = $1::varchar
		WHERE id = $2::varchar`,
		pass, userId)
	if err != nil {
		return
	}

	rows.Close()

	reply.Status = 202
	// SendMsg(Message{Method: "ChangePassword", Name: t.Name, Password: t.Password, Email: mail})
	return
}

func DisableAccount(args PlugRequest, reply *PlugRequest, userId string) (err error) {
	reply.Status = 404
	if userId == "" {
		err = errors.New("User id needed for desactivation")
		return
	}

	reply.Status = 500
	rows, err := db.Query(
		`UPDATE users
		SET activated = false
		WHERE id = $1::varchar`,
		userId)

	if err != nil {
		return
	}
	rows.Close()
	reply.Status = 202

	return
}

func Delete(args PlugRequest, reply *PlugRequest, userId string) (err error) {
	if len(userId) == 0 {
		reply.Status = 400
		err = errors.New("User id needed for deletion")
		log.Warn(err)
		return
	}

	reply.Status = 500
	rows, err := db.Query("DELETE FROM users WHERE id = $1::varchar", userId)
	if err != nil {
		log.Error(err)
		return
	}
	rows.Close()
	// SendMsg(Message{Method: "Delete", Email: mail})

	reply.Status = 202
	return
}

func ListCall(args PlugRequest, reply *PlugRequest, mail string) error {
	users, err := GetUsers()
	if err != nil {
		return err
	}

	rsp, err := json.Marshal(users)
	reply.Body = string(rsp)
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json; charset=UTF-8"
	if err == nil {
		reply.Status = 200
	} else {
		reply.Status = 400
	}
	return nil
}

var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string) error
}{
	{`^\/api\/users\/(?P<id>[^\/]+)\/disable\/{0,1}$`, "POST", DisableAccount},
	{`^\/api\/users\/{0,1}$`, "GET", ListCall},
	{`^\/api\/users\/{0,1}$`, "POST", Add},
	{`^\/api\/users\/(?P<id>[^\/]+)\/{0,1}$`, "DELETE", Delete},
	{`^\/api\/users\/(?P<id>[^\/]+)\/{0,1}$`, "PUT", UpdatePassword},
	{`^\/api\/users\/(?P<id>[^\/]+)\/{0,1}$`, "GET", GetUser},
}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			fmt.Fprintf(os.Stderr, ">> %s\n", val.Url)
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				err := val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
				if err != nil {
					log.Error(err)
				}
			} else {
				err := val.f(args, reply, "")

				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	return nil
}

type Queue struct {
	Name string
}

func getUserFromEmailPassword(email, password string) (*UserInfo, string, error) {
	log.Debug("getUserFromEmailPassword")
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

	var user UserInfo
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

func (api) GetUser(arg struct{ UserId string }, res *struct {
	Success      bool
	ErrorMessage string
	User         UserInfo
}) error {
	rows, err := db.Query(
		`SELECT id, activated,
		email,
		first_name, last_name,
		is_admin, sam, windows_password
		FROM users
		WHERE id = $1::varchar`,
		arg.UserId)

	if err != nil {
		return err
	}

	defer rows.Close()
	if !rows.Next() {
		res.ErrorMessage = "User not found"
		return nil
	}

	err = rows.Scan(
		&res.User.Id, &res.User.Activated,
		&res.User.Email, &res.User.FirstName,
		&res.User.LastName, &res.User.IsAdmin,
		&res.User.Sam, &res.User.WindowsPassword)
	if err != nil {
		return err
	}
	res.Success = true
	return nil
}

func (api) AuthenticateUser(info struct {
	Username string
	Password string
}, res *struct {
	Success      bool
	ErrorMessage string
	User         UserInfo
}) error {
	user, message, err := getUserFromEmailPassword(info.Username, info.Password)

	if err != nil {
		return err
	}

	if user != nil {
		res.Success = true
		res.User = *user
		return nil
	}

	res.ErrorMessage = message
	return nil
}

func ListenToQueue() {
	conn, err := amqp.Dial(conf.QueueUri)
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		"users_topic", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare an exchange")
	_, err = ch.QueueDeclare(
		"users", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare an queue")
	err = ch.QueueBind(
		"users",       // queue name
		"*.users",     // routing key
		"users_topic", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
	responses, err := ch.Consume(
		"users", // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for d := range responses {
			HandleReturns(d.Body)
		}
	}()
	log.Println("Waiting for responses of other plugins")
	defer ch.Close()
	defer conn.Close()
	<-forever
}

func HandleReturns(ret []byte) {
	var Msg ReturnMsg
	err := json.Unmarshal(ret, &Msg)
	if err != nil {
		log.Println(err)
	}
	if Msg.Err == "" {
		log.Println("Request:", Msg.Method, "Successfully completed by plugin", Msg.Plugin)
	} else {
		if Msg.Method == "Add" {
			log.Println("Request:", Msg.Method, "Didn't complete by plugin", Msg.Plugin, ", now reversing process")
			Delete(PlugRequest{}, &PlugRequest{}, Msg.Email)
		} else {
			log.Println("Request:", Msg.Method, "Didn't complete by plugin", Msg.Plugin)
		}
	}
}

func (api) Plug(args interface{}, reply *bool) error {
	go ListenToQueue()
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
		log.Info("[Users] users table already set up\n")
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
		log.Errorf("[Users] Unable to create users table: %s\n", err)
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
		log.Errorf("[Users] Unable to create the default user: %s\n", err)
		return err
	}

	return nil
}

func handleRPCGetUsers() ([]byte, error) {
	users, err := GetUsers()
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	res["users"] = users

	return json.Marshal(res)
}

func main() {
	var err error

	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	srv = pie.NewProvider()

	if err = srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	initConf()

	go rpcListen(conf.QueueUri, func(req map[string]interface{}) (int, []byte, error) {
		if req["action"] == "get_users" {
			res, err := handleRPCGetUsers()
			if err != nil {
				return 500, nil, err
			}
			return 200, res, nil
		}
		return 400, []byte(`{"error": "invalid action"}`), nil
	})

	db, err = sql.Open("postgres", conf.DatabaseUri)
	if err != nil {
		log.Fatalf("Cannot connect to Postgres Database: %s", err)
	}

	err = setupDb()
	if err != nil {
		log.Fatalf("[Users] unable to setup users table: %s", err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)

}
