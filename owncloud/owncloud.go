/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

var (
	apiUrl string
	client *http.Client
)

type Ocs struct {
	XMLName xml.Name `xml:"ocs"`
	Meta    OcsMeta
}
type OcsMeta struct {
	XMLName    xml.Name `xml:"meta"`
	Status     string   `xml:"status"`
	StatusCode int      `xml:"statuscode"`
	Message    string   `xml:"message"`
}

func init() {
	client = &http.Client{
		// this code should be only on debug mode
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
}

func Configure() {
	apiUrl = fmt.Sprintf("%s://%s/ocs/v1.php/cloud", conf.protocol, conf.hostname)
}

func Create(username, password string) (Ocs, error) {
	return ocsRequest("POST", apiUrl+"/users", url.Values{
		"userid":   {username},
		"password": {password},
	})
}

func Delete(username string) (Ocs, error) {
	log.Println(apiUrl + "/users/" + username)
	return ocsRequest("DELETE", apiUrl+"/users/"+username, nil)
}

// Allows to edit attributes related to a user.
// The key could be one of these values :
// email, display, password or quota.
// Only admins can edit the quota value.
func Edit(username string, key string, value string) (Ocs, error) {
	return ocsRequest("PUT", apiUrl+"/users/"+username, url.Values{
		"key":   {key},
		"value": {value},
	})
}

func ocsRequest(method, url string, data url.Values) (Ocs, error) {
	var o Ocs
	// do a new request to the api
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return o, err
	}
	req.SetBasicAuth(conf.adminLogin, conf.adminPassword)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	rsp, err := client.Do(req)
	if err != nil {
		return o, err
	}
	// verify the http status code
	if rsp.StatusCode >= 300 {
		return o, errors.New(fmt.Sprintf("HTTP error: %s", rsp.Status))
	}
	// read the response and parse it
	b, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return o, err
	}
	log.Println(string(b))
	err = xml.Unmarshal(b, &o)
	if err != nil {
		return o, err
	}
	// 100 is the successful status code of owncloud
	if o.Meta.StatusCode != 100 {
		err = errors.New(fmt.Sprintf("Owncloud error %d: %s", o.Meta.StatusCode, o.Meta.Status))
	}
	return o, err
}
