package main

import (
	log "github.com/Sirupsen/logrus"
	"io"
	"net/rpc"
)

type plugin struct {
	name   string
	client *rpc.Client
}

func (p plugin) Plug() {
	var reply bool
	err := p.client.Call(p.name+".Plug", nil, &reply)
	if err != nil {
		log.Error("Error while calling Plug: ", p.name, err)
	} else {
		log.Info("Plugin " + p.name + " plugged")
	}
}

func (p plugin) Check() bool {
	reply := false
	err := p.client.Call(p.name+".Check", nil, &reply)
	if err != nil {
		log.Error("Error while calling Check: ", p.name, err)
	} else {
		log.Info("Plugin " + p.name + " checked")
	}
	return reply
}

func (p plugin) Unplug() {
	var reply bool
	err := p.client.Call(p.name+".Unplug", nil, &reply)
	if err != nil && err != io.ErrUnexpectedEOF {
		log.Error("Error while calling Unplug: ", p.name, err)
	}
	p.client.Close()
	log.Println("Plugin " + p.name + " unplugged")
}
