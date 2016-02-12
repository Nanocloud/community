package router

import (
	"github.com/Nanocloud/community/nanocloud/models/users"
)

type Request struct {
	Query  map[string][]string
	Body   []byte
	User   *users.User
	Params map[string]string
}
