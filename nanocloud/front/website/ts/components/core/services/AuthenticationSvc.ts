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

/// <reference path="../../../../../typings/tsd.d.ts" />
import * as $ from "jquery";

"use strict";

export class AuthenticationSvc {

	static $inject = [
		"$http"
	];
	constructor(
		private $http: angular.IHttpService
	) {
	}

	login(credentials: any): angular.IPromise<any> {
		let appKey = "9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae";
		let appSecret = "9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341";

		let basic = btoa(appKey + ":" + appSecret);
		console.log(basic);
		return this.$http.post("/oauth/token", JSON.stringify({
			"username": credentials.email,
			"password": credentials.password,
			"grant_type": "password"
		}), {
			headers: { "Authorization": "Basic " + basic }
		});
	}

	logout(): angular.IPromise<any> {
		return this.$http.get("/logout");
	}

}

angular.module("haptic.core").service("AuthenticationSvc", AuthenticationSvc);
