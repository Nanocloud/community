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

import "testing"

func TestSingle(t *testing.T) {
	key := "NANOCLOUD"
	value := "nanocloud"

	Set(key, value, false)

	config := Get(false, key)
	v, ok := config[key]
	if !ok || v != value {
		t.Errorf("Value should be set")
	}

	Unset(key)

	config = Get(false, key)
	_, ok = config[key]
	if ok {
		t.Errorf("Value should have been deleted")
	}
}

func TestMultiple(t *testing.T) {
	Set("LAST", "last", false)

	keys := []string{"NANOCLOUD", "FOO", "BAR", "BAZ"}
	values := []string{"nanocloud", "foo", "bar", "baz"}

	for i, key := range keys {
		Set(key, values[i], false)
	}

	config := Get(false, keys...)

	for i, key := range keys {
		v, ok := config[key]
		if !ok || v != values[i] {
			t.Errorf("Value should be set")
			return
		}
	}

	Unset(keys...)

	config = Get(false, keys...)

	for _, key := range keys {
		_, ok := config[key]
		if ok {
			t.Errorf("Value should have been deleted")
			return
		}
	}

	last := Get(false, "LAST")["LAST"]
	if last != "last" {
		t.Errorf("Last value should be set")
	}

	Unset("LAST")

	_, ok := Get(false, "LAST")["last"]
	if ok {
		t.Errorf("Last value should have been deleted")
	}
}
