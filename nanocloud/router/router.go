package router

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Nanocloud/nano"
)

type hash map[string]interface{}
type handler func(*Request) (*Response, error)
type reqHandler struct {
	handlers *[]handler
	pattern  string
}

var module nano.Module
var kHandlers map[string][]*reqHandler

const (
	URLPrefix = "/api"
)

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

func addHandler(method string, pattern string, handlers *[]handler) {
	if kHandlers == nil {
		kHandlers = make(map[string][]*reqHandler)
	}

	h, exists := kHandlers[method]
	if !exists {
		kHandlers[method] = []*reqHandler{
			&reqHandler{
				pattern:  pattern,
				handlers: handlers,
			},
		}
	}

	kHandlers[method] = append(h, &reqHandler{
		pattern:  pattern,
		handlers: handlers,
	})
}

func Get(pattern string, handlers ...handler) {
	addHandler("GET", pattern, &handlers)
}

func Post(pattern string, handlers ...handler) {
	addHandler("POST", pattern, &handlers)
}

func Delete(pattern string, handlers ...handler) {
	addHandler("DELETE", pattern, &handlers)
}

func Put(pattern string, handlers ...handler) {
	addHandler("PUT", pattern, &handlers)
}

func Patch(pattern string, handlers ...handler) {
	addHandler("PATCH", pattern, &handlers)
}

func patternMatch(pattern, path string) (map[string]string, bool) {
	if path == "/" {
		if pattern == "/" {
			return nil, true
		}
		return nil, false
	}

	p := strings.Split(pattern, "/")
	r := strings.Split(path, "/")

	if len(p) != len(r) {
		return nil, false
	}

	m := make(map[string]string)

	for i := 0; i < len(p); i++ {
		if strings.HasPrefix(p[i], ":") {
			k := strings.TrimPrefix(p[i], ":")
			m[k] = r[i]
		} else if p[i] != r[i] {
			return nil, false
		}
	}

	return m, true
}

func handleRequest(path string, body []byte, res http.ResponseWriter, req *http.Request) *Response {
	u, err := url.Parse(req.URL.Path)
	if err != nil {
		return JSONResponse(500, hash{
			"error": err.Error(),
		})
	}

	query, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return JSONResponse(500, hash{
			"error": err.Error(),
		})
	}

	r := Request{
		Query: query,
		Body:  body,

		request:  req,
		response: res,
	}

	handlers, exists := kHandlers[req.Method]
	if exists {
		for i := 0; i < len(handlers); i++ {
			params, ok := patternMatch(handlers[i].pattern, path)
			if ok {
				r.Params = params
				var err error
				var response *Response

				for _, handler := range *handlers[i].handlers {
					response, err = handler(&r)
					if err != nil {
						return JSONResponse(500, hash{
							"error": err.Error(),
						})
					}
					if response != nil {
						return response
					}
				}

				return nil
			}
		}
	}

	contentType := ""
	rContentType, exists := req.Header["Content-Type"]
	if exists {
		if len(rContentType) != 1 {
			return JSONResponse(400, hash{
				"error": "invalid mutiple content-type",
			})
		}

		contentType = rContentType[0]
	}

	var nanoUser *nano.User

	if r.User != nil {
		nanoUser = &nano.User{
			Id:              r.User.Id,
			Email:           r.User.Email,
			Activated:       r.User.Activated,
			IsAdmin:         r.User.IsAdmin,
			FirstName:       r.User.FirstName,
			LastName:        r.User.LastName,
			Sam:             r.User.Sam,
			WindowsPassword: r.User.WindowsPassword,
		}
	}

	/* send request to RPC modules */
	response, err := module.Request(
		req.Method,
		path+"?"+req.URL.RawQuery,
		contentType,
		body,
		nanoUser,
	)

	if err != nil {
		return JSONResponse(500, hash{
			"error": err.Error(),
		})
	}

	return &Response{
		StatusCode:  response.StatusCode,
		ContentType: response.ContentType,
		Body:        response.Body,
	}
	/*
		return JSONResponse(404, hash{
			"error": "Not Found",
		})
	*/
}

func ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Cashe-Control", "no-store")
	res.Header().Set("Connection", "keep-alive")
	res.Header().Set("Expires", "Sat, 01 Jan 2000 00:00:00 GMT")
	res.Header().Set("Pragma", "no-cache")

	path := strings.TrimPrefix(req.URL.Path, URLPrefix)
	path = strings.TrimRight(path, "/")

	if len(path) == 0 {
		path = "/"
	}

	var body []byte = nil

	if req.Method == "POST" || req.Method == "PUT" || req.Method == "PATCH" {
		var err error
		body, err = ioutil.ReadAll(req.Body)
		if err != nil {
			replyError(res, 500, "unable to read request body")
			return
		}
	}

	response := handleRequest(
		path,
		body, res, req,
	)

	if response != nil {
		res.Header().Set("Content-Type", response.ContentType)
		res.WriteHeader(response.StatusCode)
		res.Write(response.Body)
	}
}

func init() {
	module = nano.RegisterModule("router")
	go module.Listen()
}
