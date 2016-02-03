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

export interface IApplication {
	alias: string;
	collection_name: string;
	Application: string;
	display_name: string;
	file_path: string;
}

export class ApplicationsSvc {

	static $inject = [
		"$http",
		"$mdToast"
	];
	constructor(
		private $http: angular.IHttpService,
		private $mdToast: angular.material.IToastService
	) {

	}

	getAll(): angular.IPromise<IApplication[]> {
		return this.$http.get("/api/apps")
			.then(
				(res: angular.IHttpPromiseCallbackArg<IApplication[]>) => {
					let apps = res.data || [];
					for (let app of apps) {
						app.Application = app.display_name;
					}
					return apps;
				},
				() => []
			);
	}

	getApplicationForUser(): angular.IPromise<IApplication[]> {
		return this.$http.get("/api/apps/me")
			.then(
				(res: angular.IHttpPromiseCallbackArg<IApplication[]>) => {
					let apps = res.data || [];
					for (let app of apps) {
						app.Application = app.display_name;
					}
					return apps;
				},
				() => []
			);
	}

	unpublish(application: IApplication): angular.IPromise<any> {
		return this.$http.delete("/api/apps/" + application.alias);
	}

}

angular.module("haptic.applications").service("ApplicationsSvc", ApplicationsSvc);
