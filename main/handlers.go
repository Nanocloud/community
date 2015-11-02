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

type test_struct struct {
	Test string
}

type Foo struct {
	g   int
	str string
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {

	str, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	t := make(map[string]string)
	t["Body"] = string(str)
	vars := mux.Vars(r)
	for i, val := range vars {
		t[i] = val
	}
	rou := html.EscapeString(r.URL.Path)
	t["path"] = rou
	//var reply Foo
	var reply map[string]string
	l := strings.Index(rou[1:], "/")
	if l == -1 {
		l = len(rou[1:])
	}
	for _, val := range plugins {
		if val.name == rou[1:l+1] {
			if val, ok := plugins["plugins/running/"+val.name]; ok {
				plugins["plugins/running/"+val.name].client.Call(plugins["plugins/running/"+val.name].name+".Receive", t, &reply)
				/*	if err != nil {
						log.Println("Error calling plugin")
						log.Println(err)
					}
				} else*/
			}
			if val, ok := plugins["plugins/staging/"+val.name]; ok {

				plugins["plugins/staging/"+val.name].client.Call(plugins["plugins/staging/"+val.name].name+".Receive", t, &reply)
				/*	if err != nil {
					log.Println("Error calling plugin")
					log.Println(err)
				*/
			}
			if reply["error"] != "404" {
				for i, val := range reply {
					w.Write([]byte(i))
					w.Write([]byte("        "))
					w.Write([]byte(val))
				}
				return
			} else {

				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 Not Found"))
				return
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
