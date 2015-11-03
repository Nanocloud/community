package main

import (
	"log"
	"net/rpc/jsonrpc"
	"os"
	"time"

	"github.com/natefinch/pie"
)

var (
	name = "plugin1" // the name should be exactly the same as the executable filename
	srv  pie.Server
	done = make(chan bool)
)

func main() {

	srv = pie.NewProvider()

	if err := srv.RegisterName(name, api{}); err != nil {
		log.Fatalf("Failed to register %s: %s", name, err)
	}

	srv.ServeCodec(jsonrpc.NewServerCodec)
}

type api struct{}

func (api) WritePage(args interface{}, reply *string) error {
	*reply = "this is written by the plugin"
	return nil
}

type Foo struct {
	g   int
	str string
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

func (api) Receive(args map[string]string, reply *map[string]string) error {
	/*	m := Foo{}
		m.g = 3
		m.str = "ff"
		*reply = m*/
	log.Println(args)
	*reply = make(map[string]string)
	if args["path"] == "/plugin1/getid" {
		(*reply)["id"] = "59"
		SetStatusOk(reply)
		return nil
	}
	if args["path"] == "/plugin1/getname" {
		(*reply)["name"] = "joe"
		SetStatusOk(reply)
		return nil
	}
	if args["path"] == "/plugin1" {
		(*reply)["available"] = "getname getid"
		SetStatusOk(reply)
		return nil
	}
	SetPageNotFound(reply)
	return nil
}

func (api) Plug(args interface{}, reply *bool) error {
	go launch()
	*reply = true
	return nil
}

func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

func (api) Unplug(args interface{}, reply *bool) error {
	defer os.Exit(0)
	// cleanup code here
	*reply = true
	return nil
}

func launch() {
	tck := time.NewTicker(time.Second)
	for {
		<-tck.C
		log.Println(name)
	}
}
