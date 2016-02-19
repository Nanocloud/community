package router

import (
	"net/http"

	"github.com/Nanocloud/community/nanocloud/models/users"
)

type Request struct {
	Query  map[string][]string
	Body   []byte
	User   *users.User
	Params map[string]string

	response http.ResponseWriter
	request  *http.Request
}

func (r *Request) Request() *http.Request {
	return r.request
}

func (r *Request) Response() http.ResponseWriter {
	return r.response
}
