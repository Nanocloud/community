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

package main

import (
	log "github.com/Sirupsen/logrus"
	"io"
	"net/rpc"
)

type plugin struct {
	name   string
	client *rpc.Client
}

func (p plugin) Plug() {
	var reply bool
	err := p.client.Call(p.name+".Plug", nil, &reply)
	if err != nil {
		log.Error("Error while calling Plug: ", p.name, err)
	} else {
		log.Info("Plugin " + p.name + " plugged")
	}
}

func (p plugin) Check() bool {
	reply := false
	err := p.client.Call(p.name+".Check", nil, &reply)
	if err != nil {
		log.Error("Error while calling Check: ", p.name, err)
	} else {
		log.Info("Plugin " + p.name + " checked")
	}
	return reply
}

func (p plugin) Unplug() {
	var reply bool
	err := p.client.Call(p.name+".Unplug", nil, &reply)
	if err != nil && err != io.ErrUnexpectedEOF {
		log.Error("Error while calling Unplug: ", p.name, err)
	}
	p.client.Close()
	log.Println("Plugin " + p.name + " unplugged")
}
