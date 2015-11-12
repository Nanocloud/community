package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//	"io"
	"log"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/natefinch/pie"
	"github.com/streadway/amqp"
)

var (
	name = "users"
	srv  pie.Server
)

type api struct{}

var UserDb *bolt.DB

type UserInfo struct {
	Name     string
	Email    string
	Password string
}

type Message struct {
	Method    string
	Name      string
	Email     string
	Activated bool
	Sam       string
	Password  string
}

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
}

func Configure() error {

	var err error
	err = nil

	UserDb, err = bolt.Open(conf.ConnectionString, 0777, nil)
	if err != nil {
		return err
	}

	err = UserDb.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(conf.DatabaseName))
		return nil
	})
	return err
}

func GetList(users *[]UserInfo) error {
	var user UserInfo

	e := UserDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			user = UserInfo{}
			json.Unmarshal(value, &user)
			*users = append(*users, user)
		}

		return nil
	})

	if e != nil {
		return e
	}

	return nil
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

var ch *amqp.Channel
var q amqp.Queue
var conn *amqp.Connection

func SendMsg(msg Message) {

	body, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

}

func Add(args PlugRequest) error {
	var msg Message
	var t UserInfo
	err := json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
	}
	msg.Method = "Add"
	msg.Name = t.Name
	msg.Email = t.Email
	msg.Password = t.Password
	SendMsg(msg)
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		bucket.Put([]byte(t.Name), []byte(args.Body))

		return err
	})

	return nil
}

func AddCall(args PlugRequest, reply *PlugRequest) {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	Add(args)
	if err != nil {
		log.Println(err)
	}
	defer UserDb.Close()
}

func ListCall(reply *PlugRequest) {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	defer UserDb.Close()
	var users []UserInfo
	GetList(&users)
	rsp, err := json.Marshal(users)
	reply.Body = string(rsp)

}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {

	if strings.Index(args.Url, "/users/add") == 0 {
		AddCall(args, reply)
	}
	if strings.Index(args.Url, "/users/list") == 0 {
		ListCall(reply)
	}

	return nil
}

func (api) Plug(args interface{}, reply *bool) error {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err = ch.QueueDeclare(
		"users", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	//go launch()
	*reply = true
	return nil
}

func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Unplug(args interface{}, reply *bool) error {
	ch.Close()
	conn.Close()
	defer os.Exit(0)
	*reply = true
	return nil
}

func main() {

	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)
}
