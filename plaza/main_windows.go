// +build windows

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
	"log"
	"os"

	"github.com/Nanocloud/community/plaza/windows/service"
	"github.com/Sirupsen/logrus"
)

const debug = false

func main() {
	if len(os.Args) < 2 || os.Args[1] != "service" {
		logrus.Info("(re)Installing service")
		err := service.InstallItSelf()
		if err != nil {
			logrus.Error(err)
		}
		return
	}

	if debug {
		out, err := os.OpenFile(`C:\plaza.log`, os.O_WRONLY|os.O_CREATE, 0644)
		if err == nil {
			defer out.Close()
			log.SetOutput(out)
			logrus.SetOutput(out)
		}
	}

	err := service.Run()
	if err != nil {
		logrus.Error(err)
	}
}
