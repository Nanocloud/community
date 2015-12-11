package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/natefinch/pie"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"
	"os"
	"regexp"
)

var (
	name = "apps"
	srv  pie.Server
)

// Structure used for exchanges with the core, faking a request/responsewriter
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

// Plugin Structure
type api struct{}

type GetApplicationsListReply struct {
	Applications []Connection
}

// Set return codes and content type of response, and list apps
func getList(args PlugRequest, reply *PlugRequest, name string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json; charset=UTF-8"
	reply.Status = 200
	connections := listApplications(reply)
	rsp, err := json.Marshal(connections)
	if err != nil {
		reply.Status = 500
		log.Error("Marshalling of connections for all users failed: ", err)
	}
	reply.Body = string(rsp)

}

// Get a list of apps accessible by the user owning the SAMAccount sam
func getListForCurrentUser(args PlugRequest, reply *PlugRequest, sam string) {

	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json; charset=UTF-8"
	reply.Status = 200
	connections := listApplicationsForSamAccount(sam, reply)

	rsp, err := json.Marshal(connections)
	if err != nil {
		reply.Status = 500
		log.WithFields(log.Fields{
			"SAMAccount": sam,
		}).Error("Marshalling of connections for current user failed: ", err)

	}
	reply.Body = string(rsp)
}

// Make an application unusable
func unpublishApplication(args PlugRequest, reply *PlugRequest, name string) {
	reply.HeadVals = make(map[string]string, 1)
	reply.HeadVals["Content-Type"] = "application/json; charset=UTF-8"
	reply.Status = 200
	if name != "" {
		unpublishApp(name)
	} else {
		reply.Status = 500
		log.Error("No Application name to unpublish")
	}
}

// slice of available handler functions
var tab = []struct {
	Url    string
	Method string
	f      func(PlugRequest, *PlugRequest, string)
}{
	{`^\/api\/apps\/{0,1}$`, "GET", getList},
	{`^\/api\/apps\/(?P<id>[^\/]+)\/{0,1}$`, "DELETE", unpublishApplication},
	{`^\/api\/apps\/(?P<id>[^\/]+)\/{0,1}$`, "GET", getListForCurrentUser},
}

// Will receive all http requests starting by /api/history from the
// core and chose the correct handler function
func (api) Receive(args PlugRequest, reply *PlugRequest) error {
	for _, val := range tab {
		re := regexp.MustCompile(val.Url)
		match := re.MatchString(args.Url)
		if val.Method == args.Method && match {
			if len(re.FindStringSubmatch(args.Url)) == 2 {
				val.f(args, reply, re.FindStringSubmatch(args.Url)[1])
			} else {
				val.f(args, reply, "")
			}
		}
	}
	return nil
}

// Plug the plugin to the core
func (api) Plug(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

// Will contain various verifications for the plugin. If the core can
// call the function and receives "true" in the reply, it means the
// plugin is functionning correctly
func (api) Check(args interface{}, reply *bool) error {
	*reply = true
	return nil
}

// Unplug the plugin from the core
func (api) Unplug(args interface{}, reply *bool) error {
	defer os.Exit(0)
	*reply = true
	return nil
}

func main() {
	var err error

	log.SetOutput(os.Stderr)
	log.SetLevel(log.DebugLevel)

	srv = pie.NewProvider()

	if err = srv.RegisterName(name, api{}); err != nil {
		log.Fatal("Failed to register %s: %s", name, err)
	}

	initConf()

	srv.ServeCodec(jsonrpc.NewServerCodec)
}
