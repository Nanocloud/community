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
	"regexp"
	//	"strings"

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
	Method   string
	UrlVars  map[string]string
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
	return e
}

func GetUser(args PlugRequest, reply *PlugRequest, mail string) error {
	if mail == "" {
		return errors.New(fmt.Sprintf("Email needed to retrieve account informations"))
	}
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
		return err
	}
	defer UserDb.Close()
	e := UserDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}

		userJson := bucket.Get([]byte(mail))
		//json.Unmarshal(userJson, &user)
		reply.Body = string(userJson)

		return nil
	})
	return e

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

func Add(args PlugRequest, reply *PlugRequest, mail string) error {
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

func ModifyPassword(args PlugRequest, reply *PlugRequest, mail string) error {
	if mail == "" {
		return errors.New(fmt.Sprintf("Email needed to modify account"))
	}
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
			if string(key) == mail {
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

func DisableAccount(args PlugRequest, reply *PlugRequest, mail string) error {
	if mail == "" {
		return errors.New(fmt.Sprintf("Email needed for desactivation"))
	}
	initConf()
	err := Configure()
	if err != nil {
		return err
	}
	defer UserDb.Close()
	var rec UserInfo
	if err != nil {
		return err
	}
	SendMsg(Message{Method: "DisableAccount", Email: mail})
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			if string(key) == mail {
				json.Unmarshal(value, &rec)
				rec.Activated = "false"
				jsonUser, _ := json.Marshal(rec)
				bucket.Put([]byte(rec.Email), jsonUser)
				break
			}
		}

		return err
	})
	return nil
}

func Delete(args PlugRequest, reply *PlugRequest, mail string) error {
	if mail == "" {
		return errors.New(fmt.Sprintf("Email needed for deletion"))
	}
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}

	defer UserDb.Close()
	SendMsg(Message{Method: "Delete", Email: mail})
	UserDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		bucket.Delete([]byte(mail))

		return err
	})

	return nil
}

func ListCall(args PlugRequest, reply *PlugRequest, mail string) error {
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

var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string) error
}{
	{`^\/users\/(?P<id>[^\/]+)\/disable\/{0,1}$`, "POST", DisableAccount},
	{`^\/users\/{0,1}$`, "GET", ListCall},
	{`^\/users\/{0,1}$`, "POST", Add},
	{`^\/users\/(?P<id>[^\/]+)\/{0,1}$`, "DELETE", Delete},
	{`^\/users\/(?P<id>[^\/]+)\/{0,1}$`, "PUT", ModifyPassword},
	{`^\/users\/(?P<id>[^\/]+)\/{0,1}$`, "GET", GetUser},
}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				err := val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
				if err != nil {
					log.Println(err)
				}
			} else {
				err := val.f(args, reply, "")

				if err != nil {
					log.Println(err)
				}
			}
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
