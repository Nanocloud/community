package main

import (
	"fmt"
	//"github.com/gorilla/mux"
	//"html"
	"html"
	//"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

type PlugRequest struct {
	Body     string
	Header   http.Header
	Form     url.Values
	PostForm url.Values
	Url      string
}

func GetRequestInfos(w http.ResponseWriter, r *http.Request, t *PlugRequest) {
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	t.Body = string(str)
	t.Header = r.Header
	err = r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	t.Form = r.Form
	t.PostForm = r.PostForm
	t.Url = html.EscapeString(r.URL.Path)

}

func CheckTrailingSlash(t PlugRequest) int {
	l := strings.Index(t.Url[1:], "/")
	if l == -1 {
		l = len(t.Url[1:])
	}
	return l
}

func WriteError(w http.ResponseWriter, reply map[string]string) {
	w.Write([]byte(reply["statuscode"]))
	w.Write([]byte(" " + reply["errormsg"]))
}

func WriteAnswer(w http.ResponseWriter, reply PlugRequest) {
	//	w.Header["Access-Control-Allow-Origin"] = "*"
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cashe-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Expires", "Sat, 01 Jan 2000 00:00:00 GMT")
	w.Header().Set("Pragma", "no-cache")
	w.Write([]byte(reply.Body))
}

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	var args PlugRequest
	GetRequestInfos(w, r, &args)
	l := CheckTrailingSlash(args)
	var rep PlugRequest
	for _, val := range plugins {
		if val.name == args.Url[1:l+1] {
			// TODO Reunite these 2 cases in 1
			if val, ok := plugins[conf.RunDir+val.name]; ok {
				err := plugins[conf.RunDir+val.name].client.Call(plugins[conf.RunDir+val.name].name+".Receive", args, &rep)
				if err != nil {
					log.Println(err)
				}
			}
			if val, ok := plugins[conf.StagDir+val.name]; ok {
				err := plugins[conf.StagDir+val.name].client.Call(plugins[conf.StagDir+val.name].name+".Receive", args, &rep)
				if err != nil {
					log.Println(err)
				}
			}
			WriteAnswer(w, rep)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
