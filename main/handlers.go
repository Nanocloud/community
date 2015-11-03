package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func GetRequestInfos(w http.ResponseWriter, r *http.Request, t map[string]string) {
	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	t["Body"] = string(str)
	vars := mux.Vars(r)
	for i, val := range vars {
		t[i] = val
	}
	rou := html.EscapeString(r.URL.Path)
	t["path"] = rou
	for i, val := range r.Header {
		for _, v := range val {
			t[i] += " " + v
		}
	}
	t["method"] = r.Method
	t["host"] = r.Host
	t["requesturi"] = r.RequestURI

}

func CheckTrailingSlash(t map[string]string) int {
	l := strings.Index(t["path"][1:], "/")
	if l == -1 {
		l = len(t["path"][1:])
	}
	return l
}

func WriteError(w http.ResponseWriter, reply map[string]string) {
	w.Write([]byte(reply["statuscode"]))
	w.Write([]byte(" " + reply["errormsg"]))
}

func WriteAnswer(w http.ResponseWriter, reply map[string]string) {
	for i, val := range reply {
		if i != "statuscode" && i != "errormsg" && i != "errordesc" {
			w.Write([]byte(i))
			w.Write([]byte("        "))
			w.Write([]byte(val))
		}
	}
}

func GenericHandler(w http.ResponseWriter, r *http.Request) {
	args := make(map[string]string)
	GetRequestInfos(w, r, args)
	var reply map[string]string
	l := CheckTrailingSlash(args)
	for _, val := range plugins {
		if val.name == args["path"][1:l+1] {
			// TODO Reunite these 2 cases in 1
			if val, ok := plugins["plugins/running/"+val.name]; ok {
				plugins["plugins/running/"+val.name].client.Call(plugins["plugins/running/"+val.name].name+".Receive", args, &reply)
			}
			if val, ok := plugins["plugins/staging/"+val.name]; ok {

				plugins["plugins/staging/"+val.name].client.Call(plugins["plugins/staging/"+val.name].name+".Receive", args, &reply)
			}
			if reply["statuscode"] != "404" {
				WriteAnswer(w, reply)
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				WriteError(w, reply)
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
