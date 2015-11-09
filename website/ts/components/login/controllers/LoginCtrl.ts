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
/// <amd-dependency path="../services/AuthenticationSvc" />
import { AuthenticationSvc } from "../services/AuthenticationSvc";

"use strict";

export class LoginCtrl {

	credentials: any;

	static $inject = [
		"$location",
		"AuthenticationSvc",
		"$mdToast"
	];

	constructor(
		private $location: angular.ILocationService,
		private authSrv: AuthenticationSvc,
		private $mdToast: angular.material.IToastService
	) {
		this.credentials = {
			"email": "",
			"password": ""
		};
	}

	signIn(e: MouseEvent) {
		let user = {
			"Email": this.credentials.email
		};
		sessionStorage.setItem("user", user.Email);

		this.authSrv.authenticate(this.credentials).then(
				(response: any) => {
					// if (typeof response.data === "string") {
					// 	if (response.data.indexOf instanceof Function &&
					// 			response.data.indexOf("<body layout=\"row\" ng-controller=\"MainCtrl as mainCtrl\">") !== -1) {
					// 		this.$location.path("/admin.html");
					// 		window.location.href = "/admin.html";
					// 		return;
					// 	}
					// }
					// this.$location.path("/");
					// window.location.href = "/";
					console.log(response);
					this.$location.path("#/");
				},
				(error: any) => {
					this.$mdToast.show(
							this.$mdToast.simple()
							.content("Authentication failed: Email or Password incorrect")
							.position("top right")
							);
				}
		);
	}
}

angular.module("haptic.login").controller("LoginCtrl", LoginCtrl);
