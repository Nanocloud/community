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
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"regexp"

	"github.com/natefinch/pie"
	"github.com/streadway/amqp"

	//todo vendor this dependency
	// nan "nanocloud.com/plugins/owncloud/libnan"
)

// Create an object to be exported

var (
	name = "owncloud"
	srv  pie.Server
)
var ch *amqp.Channel
var q amqp.Queue
var conn *amqp.Connection

type CreateUserParams struct {
	Email, Password string
}
type Message struct {
	Method    string
	Name      string
	Email     string
	Activated string
	Sam       string
	Password  string
}

type api struct{}

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
	Method   string
	Status   int
	HeadVals map[string]string
}

type ReturnMsg struct {
	Method string
	Err    string
	Plugin string
	Email  string
}

func CreateUser(args PlugRequest, reply *PlugRequest, mail string) {

	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "text/html; charset=UTF-8"
	var params CreateUserParams
	err := json.Unmarshal([]byte(args.Body), &params)
	if err != nil {
		reply.Status = 400
		log.Println(err)
		return
	}
	var msg ReturnMsg
	msg = Create(params.Email, params.Password)
	if msg.Err != "" {
		reply.Status = 400
		log.Println(msg.Err)
		return
	}
	reply.Status = 201
}

func ChangePassword(args PlugRequest, reply *PlugRequest, mail string) {

	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "text/html; charset=UTF-8"
	var params CreateUserParams
	err := json.Unmarshal([]byte(args.Body), &params)
	if err != nil {
		reply.Status = 400
		log.Println(err)
		return
	}
	var msg ReturnMsg
	msg = Edit(mail, "password", params.Password)
	if msg.Err != "" {
		reply.Status = 404
		log.Println(msg.Err)
		return
	}
	reply.Status = 204
}

func DeleteUser(args PlugRequest, reply *PlugRequest, mail string) {

	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "text/html; charset=UTF-8"
	var msg ReturnMsg
	msg = Delete(mail)
	if msg.Err != "" {
		log.Println("deletion error: ", msg.Err)
		reply.Status = 404
		return
	}
	reply.Status = 204
}

var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string)
}{
	{`^\/api\/owncloud\/users\/{0,1}$`, "POST", CreateUser},
	{`^\/api\/owncloud\/users\/(?P<id>[^\/]+)\/{0,1}$`, "DELETE", DeleteUser},
	{`^\/api\/owncloud\/users\/(?P<id>[^\/]+)\/{0,1}$`, "PUT", ChangePassword},
}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	initConf()
	Configure()
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
			} else {
				val.f(args, reply, "")
			}
		}
	}
	return nil
}

/*
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
}*/

type Queue struct {
	Name string
}

func (api) Plug(args interface{}, reply *bool) error {
	*reply = true
	go LookForMsg()
	return nil
}

func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Unplug(args interface{}, reply *bool) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	ch.Close()
	conn.Close()
	defer os.Exit(0)
	*reply = true
	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func SendReturn(msg ReturnMsg) {
	Str, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
	err = ch.Publish(
		"users_topic",    // exchange
		"owncloud.users", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        Str,
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent return to users")

}

func LookForMsg() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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
	_, err = ch.QueueDeclare(
		"owncloud", // name
		true,       // durable
		false,      // delete when usused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	failOnError(err, "Failed to declare an queue")

	err = ch.QueueBind(
		"owncloud",    // queue name
		"users.*",     // routing key
		"users_topic", // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")
	msgs, err := ch.Consume(
		"owncloud", // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		var msg Message
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				log.Println(err)
			}
			HandleRequest(msg)

		}
	}()

	log.Printf(" [*] Waiting for messages from Users")
	<-forever
}

/*
func HandleError(msg ReturnMsg) {
	if err != nil {
		log.Println(err)
		SendReturn("Plugin owncloud encountered an error in the request")
	} else {
		SendReturn("Plugin owncloud successfully completed the request")
	}
}*/

func HandleRequest(msg Message) {
	initConf()
	Configure()
	var ret ReturnMsg
	if msg.Method == "Add" {
		ret = Create(msg.Email, msg.Password)
		SendReturn(ret)
	} else if msg.Method == "Delete" {
		ret = Delete(msg.Email)
		SendReturn(ret)
	} else if msg.Method == "ChangePassword" {
		ret = Edit(msg.Email, "password", msg.Password)
		SendReturn(ret)
	}

}

func main() {
	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)

}
