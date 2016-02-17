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

export interface IService {
	Ico: string;
	Name: string;
	DisplayName: string;
	Locked: boolean;
	Status: string;
	FontColor: string;
	CurrentSize: number;
	TotalSize: number;
	VM: string;
}

export class ServicesSvc {

	static $inject = [
		"$http",
		"$mdToast",
		"$mdDialog"
	];

	downloadStarted = false;
	windowsReady = false;

	constructor(
		private $http: angular.IHttpService,
		private $mdToast: angular.material.IToastService,
		private $mdDialog: angular.material.IDialogService
	) {

	}

	getAll(): angular.IPromise<IService[]> {
		return this.$http.get("/api/iaas").then(
			(res: angular.IHttpPromiseCallbackArg<IService[]>) => {
				if (!res.data) {
					return [];
				}
				for (let srv of res.data) {
					if (srv.Ico === "windows") {
						this.downloadStarted = false;
						this.windowsReady = true;
						break;
					}
				}
				return res.data;
			},
			() => []);
	}

	download(serviceName: string): angular.IPromise<any> {
		return this.$http.post("/api/iaas/" + serviceName + "/download", null)
			.then(() => {
				this.downloadStarted = true;
			});
	}

	start(service: IService): angular.IPromise<any> {
		return this.$http.post("/api/iaas/" + service.Name + "/start", null).then(
			function() {
				service.Status = "booting";
			},
			function() {
				service.Status = "available";
			});
	}

	startStopService(service: IService): void  {
		let o = this.$mdDialog.confirm()
			.parent(angular.element(document.body))
			.title("Stop service")
			.textContent("Are you sure you want to stop this service ?")
			.ok("Yes")
			.cancel("No");
		let e = this.$mdDialog
			.show(o).then(() => {
				this.stop(service);
			});
	}

	stop(service: IService): angular.IPromise<any> {
		return this.$http.post("/api/iaas/" + service.Name + "/stop", null).then(
			function() {
				service.Status = "available";
			},
			function() {
				service.Status = "running";
			});
	}

}

angular.module("haptic.services").service("ServicesSvc", ServicesSvc);
