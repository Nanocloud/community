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
	Name      string
	Email     string
	Password  string
	Activated string
	Sam       string
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

func GetUser(args PlugRequest, reply *PlugRequest) error {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
		return err
	}
	defer UserDb.Close()
	var user UserInfo
	err = json.Unmarshal([]byte(args.Body), &user)
	if err != nil {
		log.Println(err)
		return err
	}
	e := UserDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}

		userJson := bucket.Get([]byte(user.Email))
		//json.Unmarshal(userJson, &user)
		reply.Body = string(userJson)

		return nil
	})
	if e != nil {
		return e
	}
	return nil

}

type Registered struct {
	IsRegistered string
}

func IsUserRegistered(args PlugRequest, reply *PlugRequest) error {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
		return err
	}
	defer UserDb.Close()
	var user UserInfo
	err = json.Unmarshal([]byte(args.Body), &user)
	if err != nil {
		log.Println(err)
		return err
	}

	str, err := json.Marshal(Registered{IsRegistered: "false"})
	if err != nil {
		log.Println(err)
	}
	reply.Body = string(str)
	e := UserDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}

		cursor := bucket.Cursor()
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			if string(key) == user.Email {
				str, err := json.Marshal(Registered{IsRegistered: "true"})
				if err != nil {
					log.Println(err)
				}
				reply.Body = string(str)
				break
			}
		}

		return nil
	})
	return e
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type NbUser struct {
	Count int
}

func CountRegisteredUsers(args PlugRequest, reply *PlugRequest) error {
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	var nb NbUser
	nb.Count = 0
	defer UserDb.Close()

	e := UserDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}

		cursor := bucket.Cursor()
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			nb.Count++
		}
		str, err := json.Marshal(nb)
		if err != nil {
			log.Println(err)
		}
		reply.Body = string(str)
		return nil
	})

	return e
}

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
	defer ch.Close()
	defer conn.Close()

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
	SendMsg(Message{Method: "Add", Name: t.Name, Email: t.Email, Password: t.Password, Activated: t.Activated})
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		bucket.Put([]byte(t.Email), []byte(args.Body))

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
	SendMsg(Message{Method: "ChangePassword", Name: t.Name, Password: t.Password, Email: t.Email})
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
				bucket.Put([]byte(rec.Email), jsonUser)
				break
			}
		}

		return err
	})
	return nil
}

func DisableAccount(args PlugRequest, reply *PlugRequest) error {
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
	SendMsg(Message{Method: "DisableAccount", Email: t.Email, Name: t.Name})
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == t.Email {
				json.Unmarshal(value, &rec)
				rec.Activated = t.Activated
				jsonUser, _ := json.Marshal(rec)
				bucket.Put([]byte(rec.Email), jsonUser)
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
	SendMsg(Message{Method: "Delete", Email: t.Email, Name: t.Name})
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		bucket.Delete([]byte(t.Email))

		return err
	})

	return nil
}

func ListCall(args PlugRequest, reply *PlugRequest) error {
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
	return nil

}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	handlers := make(map[string]func(PlugRequest, *PlugRequest) error, 8)
	handlers["/users/add"] = Add
	handlers["/users/list"] = ListCall
	handlers["/users/delete"] = Delete
	handlers["/users/modifypassword"] = ModifyPassword
	handlers["/users/disableaccount"] = DisableAccount
	handlers["/users/getuser"] = GetUser
	handlers["/users/isuserregistered"] = IsUserRegistered
	handlers["/users/countregisteredusers"] = CountRegisteredUsers

	for i, val := range handlers {
		if strings.Index(args.Url, i) == 0 {
			val(args, reply)
			break
		}
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
			if err != nil {
				log.Println(err)
			}
			log.Println(string(d.Body))
		}
	}()
	log.Println("Waiting for responses of fake/owncloud")
	defer ch.Close()
	defer conn.Close()
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
