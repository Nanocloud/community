package main

import (
	"encoding/json"
	"github.com/Nanocloud/nano"
	"github.com/Nanocloud/oauth"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpHandler struct {
	URLPrefix string
	Module    nano.Module
}

func replyError(res http.ResponseWriter, statusCode int, message string) {
	m := make(map[string]string)

	m["error"] = message

	b, err := json.Marshal(m)
	if err != nil {
		res.Write([]byte(`{"error":"internal server error"}`))
		res.WriteHeader(500)
		return
	}

	res.Write(b)
	res.WriteHeader(statusCode)
}

func (h httpHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Cashe-Control", "no-store")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Expires", "Sat, 01 Jan 2000 00:00:00 GMT")
	res.Header().Set("Pragma", "no-cache")

	path := strings.TrimPrefix(req.URL.Path, h.URLPrefix)
	path = strings.TrimRight(path, "/")

	if len(path) == 0 {
		path = "/"
	}

	contentType := ""
	var body []byte = nil

	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
		rContentType, exists := req.Header["Content-Type"]
		if exists {
			if len(rContentType) != 1 {
				replyError(res, 400, "invalid mutiple content-type")
				return
			}

			contentType = rContentType[0]
		}

		var err error
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			replyError(res, 500, "unable to read request body")
			return
		}
	}

	user := oauth.GetUserOrFail(res, req)
	if user == nil {
		return
	}

	response, err := h.Module.Request(
		req.Method,
		path+"?"+req.URL.RawQuery,
		contentType,
		body,
		user.(*nano.User),
	)

	if err != nil {
		module.Log.Error(err)
		res.WriteHeader(500)
		return
	}

	res.Header().Set("Content-Type", response.ContentType)
	res.Write(response.Body)
	res.WriteHeader(response.StatusCode)
}
