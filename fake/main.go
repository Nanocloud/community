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

	"github.com/natefinch/pie"
	"github.com/streadway/amqp"

	//todo vendor this dependency
	// nan "nanocloud.com/plugins/fake/libnan"
)

// Create an object to be exported

var (
	name = "fake"
	srv  pie.Server
)
var ch *amqp.Channel
var q amqp.Queue
var conn *amqp.Connection

type Message struct {
	Method    string
	Name      string
	Email     string
	Activated bool
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
}

func Create() error {
	log.Println("I CREATED A FAKE USER IN A FAKE DB YAY")
	return nil
}

type del struct {
	Username string
}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {

	return nil
}

type Queue struct {
	Name string
}

func SendMyQueues() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"users", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	tmp := Queue{Name: "users.fake"}
	str, err := json.Marshal(tmp)
	log.Println("json sent to users plugin", string(str))
	if err != nil {
		log.Println(err)
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "encoding/json",
			Body:        str,
		})
	log.Printf(" [x] Sent a Q users.fake")
	failOnError(err, "Failed to publish a message")
	tmp = Queue{Name: "fake.users"}
	str, err = json.Marshal(tmp)
	if err != nil {
		log.Println(err)
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "encoding/json",
			Body:        str,
		})
	log.Printf(" [x] Sent a Q fake.users")
	failOnError(err, "Failed to publish a message")

}

func (api) Plug(args interface{}, reply *bool) error {
	*reply = true
	SendMyQueues()
	go LookForMsg()
	return nil
}

func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Unplug(args interface{}, reply *bool) error {
	SendReturn("stop")
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

func SendReturn(msg string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	q, err := ch.QueueDeclare(
		"fake.users", // name
		false,        // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")
	tmp := Queue{Name: msg}
	str, err := json.Marshal(tmp)
	log.Println("json sent to users plugin", string(str))
	if err != nil {
		log.Println(err)
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "encoding/json",
			Body:        str,
		})
	failOnError(err, "Failed to publish a message")
}

func LookForMsg() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"users.fake", // name
		false,        // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
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
			if msg.Method == "Add" {
				err := Create()
				if err != nil {
					log.Println("create error?:")
					log.Println(err)
					SendReturn("Plugin fake encounter an error in the request")
				} else {
					SendReturn("Plugin fake successfully completed the request")
				}
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
func main() {
	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)

}
