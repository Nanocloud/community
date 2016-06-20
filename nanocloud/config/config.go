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

package config

import (
	"fmt"
	"strings"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	log "github.com/Sirupsen/logrus"
)

// Return a map with the keys arguments associated with the found value if any.
// If the key doesn't exist, no field for the actual key will be present in the map.
func Get(keys ...string) map[string]string {

	rt := make(map[string]string)

	l := len(keys)
	if l == 0 {
		return rt
	}

	queryArgs := make([]string, l)
	args := make([]interface{}, l)
	for k, v := range keys {
		args[k] = v
		queryArgs[k] = fmt.Sprintf("$%d::varchar", k+1)
	}

	request := fmt.Sprintf("SELECT key, value FROM config WHERE key IN(%s)", strings.Join(queryArgs, ","))
	rows, err := db.Query(request, args...)

	if err != nil {
		log.Error(err)
		return rt
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var value string

		rows.Scan(&key, &value)
		rt[key] = value
	}

	return rt
}

// Save the configuration key. If the value exists already, it will be overwritten.
func Set(key, value string) {
	_, err := db.Exec(
		`INSERT INTO
		config (key, value)
		VALUES($1::varchar, $2::varchar)
		ON CONFLICT(key) DO UPDATE SET value = excluded.value`,
		key, value,
	)
	if err != nil {
		log.Error(err)
	}
}

// Remove the specified configuration keys.
func Unset(keys ...string) {
	l := len(keys)

	queryArgs := make([]string, l)
	args := make([]interface{}, l)
	for k, v := range keys {
		args[k] = v
		queryArgs[k] = fmt.Sprintf("$%d::varchar", k+1)
	}

	_, err := db.Exec(
		fmt.Sprintf("DELETE FROM config WHERE key IN(%s)", strings.Join(queryArgs, ",")),
		args...,
	)
	if err != nil {
		log.Error(err)
	}
}
