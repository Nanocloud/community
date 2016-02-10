package router

import (
	"encoding/json"
	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/oauth2"
	"github.com/Nanocloud/nano"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type hash map[string]interface{}
type handler func(Request) (*Response, error)
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

func handleRequest(path string, user *users.User, body []byte, res http.ResponseWriter, req *http.Request) *Response {
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
		User:  user,
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
					response, err = handler(r)
					if err != nil {
						return JSONResponse(500, hash{
							"error": err.Error(),
						})
					}
					if response != nil {
						return response
					}
				}

				return JSONResponse(500, hash{
					"error": err.Error(),
				})
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

	/* send request to RPC modules */
	response, err := module.Request(
		req.Method,
		path+"?"+req.URL.RawQuery,
		contentType,
		body,
		&nano.User{
			Id:              user.Id,
			Email:           user.Email,
			Activated:       user.Activated,
			IsAdmin:         user.IsAdmin,
			FirstName:       user.FirstName,
			LastName:        user.LastName,
			Sam:             user.Sam,
			WindowsPassword: user.WindowsPassword,
		},
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

	user := oauth2.GetUserOrFail(res, req)
	if user == nil {
		return
	}

	response := handleRequest(
		path,
		user.(*users.User),
		body, res, req,
	)

	res.Header().Set("Content-Type", response.ContentType)
	res.WriteHeader(response.StatusCode)
	res.Write(response.Body)
}

func init() {
	module = nano.RegisterModule("router")
}
