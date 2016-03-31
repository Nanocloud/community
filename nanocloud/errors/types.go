/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2016 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const unableToSerializeError = `{
	"errors": [{
		"code": "0000001",
		"title": "An unexpected error occured.",
		"detail": "The system is not able to retreive the error details.",
	}]
}`

type hash map[string]interface{}

type APIError interface {
	Send(http.ResponseWriter)
}

func sendError(w http.ResponseWriter, code int, body interface{}) {
	r, err := json.Marshal(body)
	if err != nil {
		code = http.StatusInternalServerError
		r = []byte(unableToSerializeError)
	}

	w.Header().Set("Content-Type", "application/vnd.api+json")
	w.WriteHeader(code)
	w.Write(r)
}

type apiError struct {
	code   int64
	status int // HTTP Status Code
	title  string
}

func (e *apiError) Send(w http.ResponseWriter) {
	b := hash{
		"errors": [1]hash{
			hash{
				"code":  fmt.Sprintf("%06x", e.code),
				"title": e.title,
			},
		},
	}
	sendError(w, e.status, b)
}

func (e *apiError) Detail(detail string) *detailedError {
	return &detailedError{
		err:    e,
		detail: detail,
	}
}

func (e *apiError) Error() string {
	return e.title
}

type detailedError struct {
	err    *apiError
	detail string
}

func (e *detailedError) Error() string {
	return e.err.Error() + " " + e.detail
}

func (e *detailedError) Send(w http.ResponseWriter) {
	b := hash{
		"errors": [1]hash{
			hash{
				"code":   fmt.Sprintf("%06x", e.err.code),
				"title":  e.err.title,
				"detail": e.detail,
			},
		},
	}
	sendError(w, e.err.status, b)
}
