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

import "github.com/Nanocloud/community/plaza/server/router"

/*
func connectToLogger() {
	fmt.Printf("Dialling logger... ")
	conn, err := net.DialTimeout("tcp", "127.0.0.1:8081", 10*time.Second)
	fmt.Println("done")
	if err != nil {
		fmt.Println(err)
		return
	}
	if conn != nil {
		log.SetOutput(conn)
	}
}

func main() {
	if len(os.Args) < 2 || os.Args[1] != "service" {
		log.Println("(re)Installing service")
		err := service.InstallItSelf()
		if err != nil {
			log.Println(err)
		}
		return
	}
	connectToLogger()

	err := service.Run()
	if err != nil {
		log.Println(err)
	}
}*/

func main() {
	router.Start()
}
