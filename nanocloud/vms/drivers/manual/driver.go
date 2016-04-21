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

package manual

import (
	"strings"

	"github.com/Nanocloud/community/nanocloud/connectors/db"
	"github.com/Nanocloud/community/nanocloud/vms"
	log "github.com/Sirupsen/logrus"
)

type driver struct{}

func Find(ip string) bool {

	rows, _ := db.Query(`SELECT execserver
	FROM machines WHERE execserver = $1::varchar`, ip)
	if rows.Next() {
		return true
	}
	return false
}

func (d *driver) Open(options map[string]string) (vms.VM, error) {
	v := &vm{}
	ad := options["ad"]
	servers := options["servers"]
	password := options["password"]
	user := options["user"]
	attr := vms.MachineAttributes{
		Type:     nil,
		Username: user,
		Password: password,
	}

	if ad == servers {

		if Find(ad) == false {
			attr.Ip = ad
			_, err := v.Create(attr)
			if err != nil {
				log.Error(err)
			}
		}
	} else {

		ips := strings.Split(servers, ";")
		if Find(ad) == false {
			attr.Ip = ad
			_, err := v.Create(attr)
			if err != nil {
				log.Error(err)
			}
		}

		for _, val := range ips {
			if Find(val) == false {
				attr.Ip = val
				_, err := v.Create(attr)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
	return v, nil
}
