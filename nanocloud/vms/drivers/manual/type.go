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

type machineType struct {
	id    string
	label string
	size  string
	cpu   int
	ram   int
}

func (t *machineType) Id() string {
	return t.id
}

func (t *machineType) Label() string {
	return t.label
}

var defaultType *machineType

func init() {
	defaultType = &machineType{
		id:    "default",
		label: "Default",
		size:  "60GB",
		cpu:   2,
		ram:   4096,
	}
}
