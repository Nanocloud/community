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

func Add(args PlugRequest) error {
	var t HistoryInfo
	err := json.Unmarshal([]byte(args.Body), &t)
	if err != nil {
		log.Println(err)
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

func SetStatusOk(reply *map[string]string) {
	(*reply)["statuscode"] = "200"
	(*reply)["errormsg"] = ""
	(*reply)["errordesc"] = ""
}

func SetPageNotFound(reply *map[string]string) {
	(*reply)["statuscode"] = "404"
	(*reply)["errormsg"] = "Page Not Found"
	(*reply)["errordesc"] = "The page was not found"
}

func AddCall(args PlugRequest, reply *PlugRequest) {
	log.Println("ADDCALL")
	//	ids := strings.Split(args.Url, "/")
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	Add(args)
	if err != nil {
		log.Println(err)
	}
	defer HistoryDb.Close()
}

func ListCall(reply *PlugRequest) {
	log.Println("LIST CALL")
	initConf()
	err := Configure()
	if err != nil {
		log.Println(err)
	}
	defer HistoryDb.Close()
	var histories []HistoryInfo
	GetList(&histories)
	rsp, err := json.Marshal(histories)
	reply.Body = string(rsp)

}

func (api) Receive(args PlugRequest, reply *PlugRequest) error {

	if strings.Index(args.Url, "/history/add") == 0 {
		AddCall(args, reply)
	}
	if strings.Index(args.Url, "/history/list") == 0 {
		ListCall(reply)
	}

	return nil
}

func (api) Plug(args interface{}, reply *bool) error {
	//go launch()
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
