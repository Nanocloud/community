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

package me

import (
	"encoding/json"
	"net/http"

	"github.com/Nanocloud/community/nanocloud/oauth2"
	log "github.com/Sirupsen/logrus"
)

func Get(w http.ResponseWriter, r *http.Request) {
	user := oauth2.GetUserOrFail(w, r)
	if user != nil {
		b, err := json.Marshal(map[string]interface{}{"data": user})
		if err != nil {
			log.Error(err)
			w.WriteHeader(500)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}
