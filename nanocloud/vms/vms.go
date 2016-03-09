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

package vms

import "errors"

var drivers map[string]Driver

func Register(name string, driver Driver) {
	if drivers == nil {
		drivers = make(map[string]Driver, 0)
	}
	drivers[name] = driver
}

func Open(driverName string, options map[string]string) (*VM, error) {
	driver := drivers[driverName]
	if driver != nil {
		d, _ := driver.Open(options)
		return &d, nil
	}
	return nil, errors.New("Invalid driver name")
}
