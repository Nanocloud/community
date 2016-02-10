package router

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

type Response struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

func JSONResponse(statusCode int, body interface{}) *Response {
	res := Response{
		ContentType: "application/json",
	}

	b, err := json.Marshal(body)
	if err != nil {
		res.StatusCode = 500
		log.Error(err)
		res.Body = []byte(`{"error":"Internal Server Error"}`)
		return &res
	}

	res.StatusCode = statusCode
	res.Body = b
	return &res
}
