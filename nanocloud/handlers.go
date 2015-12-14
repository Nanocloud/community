package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/nanocloud/oauth"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
	Method   string
	Status   int
	HeadVals map[string]string
}

type UserInfo struct {
	Id        string
	Activated bool
	Email     string
	FirstName string
	LastName  string
	IsAdmin   bool
}

type RPCRequest struct {
	Method string
	Path   string
	Body   []byte
}

// fill a fake http request for the plugins
func getRequestInfos(w http.ResponseWriter, r *http.Request, t *PlugRequest) {
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	t.Body = string(str)
	t.Header = r.Header
	t.Method = r.Method
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	t.Form = r.Form
	t.PostForm = r.PostForm
	t.Url = html.EscapeString(r.URL.Path)
}

// fill the http response from the plugins
func writeAnswer(w http.ResponseWriter, reply PlugRequest) {
	w.Header().Set("Content-Type", reply.HeadVals["Content-Type"])
	w.WriteHeader(reply.Status)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cashe-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Expires", "Sat, 01 Jan 2000 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Write([]byte(reply.Body))
}

func handleMeRequest(user *UserInfo, w http.ResponseWriter, r *http.Request) {
	me := make(map[string]interface{})

	me["id"] = user.Id
	me["first_name"] = user.FirstName
	me["last_name"] = user.LastName
	me["email"] = user.Email
	me["activated"] = user.Activated
	me["is_admin"] = user.IsAdmin

	rt, err := json.Marshal(me)
	if err != nil {
		log.Error(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.Write(rt)
}

// handle REST API
func genericHandler(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/oauth") {
		oauth.HandleRequest(w, r)
		return
	}

	user := oauth.GetUserOrFail(w, r)

	if user == nil {
		return
	}

	if r.URL.Path == "/api/me" {
		handleMeRequest(user.(*UserInfo), w, r)
		return
	}

	r.ParseForm()

	var args PlugRequest
	getRequestInfos(w, r, &args)
	var rep PlugRequest
	for _, val := range plugins {
		if strings.HasPrefix(args.Url, "/api/"+val.name) {
			if val, ok := plugins[conf.RunDir+val.name]; ok {
				err := plugins[conf.RunDir+val.name].client.Call(plugins[conf.RunDir+val.name].name+".Receive", args, &rep)
				if err != nil {
					log.Error(err)
				}
			}
			writeAnswer(w, rep)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

// get list of available front components
func getComponentsHandler(c *echo.Context) error {
	fis, err := ioutil.ReadDir(filepath.Join(conf.FrontDir, "ts/components"))
	if err != nil {
		log.Fatal("Unable to load the components folder. ", err)
		return c.Err()
	}
	var comps []string
	for _, f := range fis {
		comps = append(comps, f.Name())
	}
	return c.JSON(http.StatusOK, comps)
}
