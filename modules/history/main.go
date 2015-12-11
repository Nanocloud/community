package main

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"regexp"

	"github.com/boltdb/bolt"
	"github.com/natefinch/pie"
)

var (
	name = "history"
	srv  pie.Server
)

type HistoryConfig struct {
	ConnectionString string
	DatabaseName     string
}

type api struct{}

var HistoryDb *bolt.DB

type HistoryInfo struct {
	UserId       string
	ConnectionId string
	StartDate    string
	EndDate      string
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

func Configure() error {

	var err error
	err = nil

	HistoryDb, err = bolt.Open(conf.ConnectionString, 0777, nil)
	if err != nil {
		return err
	}

	err = HistoryDb.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(conf.DatabaseName))
		return nil
	})
	return err
}

func GetList(histories *[]HistoryInfo) error {
	var history HistoryInfo

	e := HistoryDb.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}

		cursor := bucket.Cursor()
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			history = HistoryInfo{}
			json.Unmarshal(value, &history)
			*histories = append(*histories, history)
		}

		return nil
	})

	if e != nil {
		return e
	}

	return nil
}

func Add(args PlugRequest, reply *PlugRequest) error {
	var t HistoryInfo
	err := json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		reply.Status = 400
		log.Error("Json Arguments are not valid: ", err)
		return err
	}

	HistoryDb.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(conf.DatabaseName))
		if bucket == nil {
			return errors.New(fmt.Sprintf("Bucket '%s' doesn't exist", conf.DatabaseName))
		}
		bucket.Put([]byte(t.ConnectionId), []byte(args.Body))

		return err
	})

	return nil
}

func AddCall(args PlugRequest, reply *PlugRequest, id string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	initConf()
	err := Configure()
	if err != nil {
		reply.Status = 500
		log.Error("Failed to connect de the Database: ", err)
		return
	}
	err = Add(args, reply)
	if err != nil {
		if reply.Status != 400 {
			reply.Status = 500
		}
		log.Error("Failed to add an entry: ", err)
		return
	}
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "text/html;charset=UTF-8"
	defer HistoryDb.Close()
	reply.Status = 202
}

func ListCall(args PlugRequest, reply *PlugRequest, id string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json;charset=UTF-8"
	initConf()
	err := Configure()
	if err != nil {
		reply.Status = 500
		log.Error("Failed to connect de the Database: ", err)
		return
	}
	defer HistoryDb.Close()
	var histories []HistoryInfo
	GetList(&histories)
	rsp, err := json.Marshal(histories)
	if err != nil {
		reply.Status = 500
		log.Error("Failed to Marshal histories: ", err)
		return
	}
	reply.Body = string(rsp)
	reply.Status = 200

}

var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string)
}{
	{`^\/api\/history\/{0,1}$`, "GET", ListCall},
	{`^\/api\/history\/{0,1}$`, "POST", AddCall},
}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
				return nil
			} else {
				val.f(args, reply, "")
				return nil
			}
		}
	}
	return nil
}

func (api) Plug(args interface{}, reply *bool) error {
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
