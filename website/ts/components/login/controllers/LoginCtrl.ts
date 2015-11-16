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
/// <amd-dependency path="../../core/services/AuthenticationSvc" />
import { AuthenticationSvc } from "../../core/services/AuthenticationSvc";

"use strict";

export class LoginCtrl {

	credentials: any;

	static $inject = [
		"$location",
		"$mdToast",
		"AuthenticationSvc"
	];

	constructor(
		private $location: angular.ILocationService,
		private $mdToast: angular.material.IToastService,
		private authSvc: AuthenticationSvc
	) {
		this.credentials = {
			"email": "",
			"password": ""
		};
	}

	signIn() {
		sessionStorage.setItem("user", this.credentials.email);
		this.authSvc
			.login(this.credentials)
			.then(
				(res: any) => {
					if (res.headers().admin && res.headers().admin === "true") {
						this.$location.path("/admin");
					} else {
						this.$location.path("/");
					}
				},
				(error: any) => {
					this.$mdToast.show(
						this.$mdToast.simple()
						.content("Authentication failed: Email or Password incorrect")
						.position("top right"));
				});
	}
}

angular.module("haptic.login").controller("LoginCtrl", LoginCtrl);
