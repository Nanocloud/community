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
/// <amd-dependency path="../../core/services/RpcSvc" />
import { RpcSvc, IRpcResponse } from "../../core/services/RpcSvc";

"use strict";

export class AuthenticationSvc {

	static $inject = [
		"$http",
		"$mdToast"
	];
	constructor(
		private $http: angular.IHttpService,
		private $mdToast: angular.material.IToastService
	) {
	}

	authenticate(credentials: any): angular.IPromise<any> {
		return this.$http.post("/login", {
			"email": credentials.email,
			"password": credentials.password
		}, {
			headers: { "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"},
			transformRequest: function(data) { return $.param(data); }
		});
	}

	private isError(res: IRpcResponse): boolean {
		if (res.error == null) {
			return false;
		}
		this.$mdToast.show(
			this.$mdToast.simple()
				.content(res.error.code === 0 ? "Internal Error" : JSON.stringify(res.error))
				.position("top right")
		);
		return true;
	}

}

angular.module("haptic.login").service("AuthenticationSvc", AuthenticationSvc);
