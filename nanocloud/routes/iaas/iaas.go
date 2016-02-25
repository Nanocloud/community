package iaas

import (
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
)

const (
	iaasAPIurl = "http://iaas-module:8080"
)

func proxy(c *echo.Context) error {
	r := c.Request()
	path := r.URL.Path

	var resp *http.Response
	var err error

	if r.Method == "GET" {
		resp, err = http.Get(iaasAPIurl + path)
	} else {
		resp, err = http.Post(iaasAPIurl+path, "", nil)
	}

	if err != nil {
		log.Error("here")
		log.Error(err)
		log.Error("there")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	contentType := ""

	ct := resp.Header["Content-Type"]
	if len(ct) > 0 {
		contentType = ct[0]
	}

	w := c.Response()
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	return nil
}

var (
	ListRunningVM = proxy
	StopVM        = proxy
	StartVM       = proxy
	DownloadVM    = proxy
)
