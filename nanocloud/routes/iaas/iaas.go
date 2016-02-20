package iaas

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Nanocloud/community/nanocloud/router"
	"github.com/labstack/gommon/log"
)

const (
	iaasAPIurl = "http://iaas-module:8080"
)

func proxy(req *router.Request) (*router.Response, error) {
	r := req.Request()
	path := r.URL.Path

	var resp *http.Response
	var err error

	if r.Method == "GET" {
		resp, err = http.Get(iaasAPIurl + path)
	} else {
		resp, err = http.Post(iaasAPIurl+path, "", nil)
	}

	if err != nil {
		log.Error(err)
		return nil, errors.New("Unable to contact Iaas API")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, errors.New("Unable to contact Iaas API")
	}

	contentType := ""

	ct := resp.Header["Content-Type"]
	if len(ct) > 0 {
		contentType = ct[0]
	}

	return &router.Response{
		StatusCode:  resp.StatusCode,
		ContentType: contentType,
		Body:        body,
	}, nil
}

var (
	ListRunningVM = proxy
	StopVM        = proxy
	StartVM       = proxy
	DownloadVM    = proxy
)
