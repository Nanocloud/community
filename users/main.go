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
	str, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	err = ch.Publish(
		"users_topic", // exchange
		"users.req",   // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "encoding/json",
			Body:        []byte(str),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent order to plugin")

}

func Add(args PlugRequest, reply *PlugRequest) error {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	defer UserDb.Close()
	var t UserInfo
	err = json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
	}
	SendMsg(Message{Method: "Add", Name: t.Name, Email: t.Email, Password: t.Password})
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

func ModifyPassword(args PlugRequest, reply *PlugRequest) error {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	defer UserDb.Close()
	var t UserInfo
	var rec UserInfo
	err = json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
	}
	SendMsg(Message{Method: "ChangePassword", Name: t.Name, Password: t.Password})
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == t.Name {
				json.Unmarshal(value, &rec)
				rec.Password = t.Password
				jsonUser, _ := json.Marshal(rec)
				bucket.Put([]byte(rec.Name), jsonUser)
				break
			}
		}

		return err
	})
	return nil
}

func Delete(args PlugRequest, reply *PlugRequest) error {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	defer UserDb.Close()
	var t UserInfo
	err = json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
	}
	SendMsg(Message{Method: "Delete", Name: t.Name})
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		bucket.Delete([]byte(t.Name))

		return err
	})

	return nil
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
		Add(args, reply)
	}
	if strings.Index(args.Url, "/users/list") == 0 {
		ListCall(reply)
	}
	if strings.Index(args.Url, "/users/delete") == 0 {
		Delete(args, reply)
	}
	if strings.Index(args.Url, "/users/modifypassword") == 0 {
		ModifyPassword(args, reply)
	}

	return nil
}

type Queue struct {
	Name string
}

func ListenToQueue() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	//defer conn.Close()

	ch, err = conn.Channel()
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
		true,    // exclusive
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
			if err != nil {
				log.Println(err)
			}
			log.Println(string(d.Body))
		}
	}()
	log.Println("Waiting for responses of fake/owncloud")
	<-forever
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
