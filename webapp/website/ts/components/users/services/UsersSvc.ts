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

"use strict";

export interface IUser {
	id?: string;
	first_name: string;
	last_name: string;
	email: string;
	profile: string;
	password?: string;
	password2?: string;
}

export class UsersSvc {

	static $inject = [
		"$http",
		"$mdToast"
	];
	constructor(
		private $http: angular.IHttpService,
		private $mdToast: angular.material.IToastService
	) {

	}

	getAll(): angular.IPromise<IUser[]> {
		return this.$http.get("/api/users")
		.then((res: any) => {
			let arr: IUser[] = [];
			for (var data of res.data.data)
			{
				let usr: IUser;
				usr = data.attributes;
				usr.id = data.id;
				arr.push(usr);
			}
			return arr;
		});
	}

	save(user: IUser): angular.IPromise<boolean> {
		return this.$http.post("/api/users", {
			data : {
        "type" : "user",
        attributes: user
      }
		})
		.then(() => true, () => false);
	}

	delete(user: IUser): angular.IPromise<any> {
		return this.$http.delete("/api/users/" + user.id);
	}

	updatePassword(user: IUser): angular.IPromise<boolean> {
		return this.$http.put("/api/users/" + user.id, {
			data: {
				password: user.password
			}
		})
		.then(() => true, () => false);
	}

}

angular.module("haptic.users").service("UsersSvc", UsersSvc);
