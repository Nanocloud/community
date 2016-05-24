/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
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

package upload

import (
	"net/http"
	"net/url"

	"github.com/Nanocloud/community/nanocloud/models/users"
	"github.com/Nanocloud/community/nanocloud/oauth2"
	"github.com/Nanocloud/community/nanocloud/utils"
	log "github.com/Sirupsen/logrus"
)

func Post(w http.ResponseWriter, r *http.Request) {
	rawuser, oauthErr := oauth2.GetUser(w, r)
	if rawuser == nil || oauthErr != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	user := rawuser.(*users.User)

	winUser, err := user.WindowsCredentials()
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	sam := winUser.Sam
	winServer := utils.Env("PLAZA_ADDRESS", "iaas-module")

	request, err := http.NewRequest(
		"POST",
		"http://"+winServer+":"+utils.Env("PLAZA_PORT", "9090")+"/upload?sam="+url.QueryEscape(sam)+"&userId="+url.QueryEscape(user.Id),
		r.Body,
	)
	if err != nil {
		log.Println("Unable de create request ", err)
	}
	request.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("Unable to send request ", err)
	}

	log.Error(resp)
	http.Error(w, "", resp.StatusCode)
	return
}

func Get(w http.ResponseWriter, r *http.Request) {
	rawuser, oauthErr := oauth2.GetUser(w, r)
	if rawuser == nil || oauthErr != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	user := rawuser.(*users.User)

	winUser, err := user.WindowsCredentials()
	if err != nil {
		http.Error(w, "", http.StatusUnauthorized)
		return
	}

	sam := winUser.Sam
	winServer := utils.Env("PLAZA_ADDRESS", "iaas-module")

	request, err := http.NewRequest(
		"GET",
		"http://"+winServer+":"+utils.Env("PLAZA_PORT", "9090")+"/upload?sam="+url.QueryEscape(sam)+"&userId="+url.QueryEscape(user.Id),
		nil,
	)
	if err != nil {
		log.Println("Unable de create request ", err)
	}

	request.Header = r.Header
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("request error", err)
	}
	if resp.StatusCode == http.StatusTeapot {
		http.Error(w, "", http.StatusSeeOther)
	}
	return
}
