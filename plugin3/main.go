package main

import (
	"log"
	"net/rpc/jsonrpc"
	"os"
	"time"

	"github.com/natefinch/pie"
)

var (
	name = "plugin3" // the name should be exactly the same as the executable filename
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
		//		log.Println(name)
		log.Println("Other plugin")
	}
}
