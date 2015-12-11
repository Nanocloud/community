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

export interface IApplication {
	Hostname: string;
	Port: string;
	Username: string;
	Password: string;
	RemoteApp: string;
	ConnectionName: string;
}

export class ApplicationsSvc {

	static $inject = [
		"RpcSvc",
		"$mdToast"
	];
	constructor(
		private rpc: RpcSvc,
		private $mdToast: angular.material.IToastService
	) {

	}

	getAll(): angular.IPromise<IApplication[]> {
		return this.rpc.call({ method: "ServiceApplications.GetList", id: 1 })
			.then((res: IRpcResponse): IApplication[] => {
				let applications: IApplication[] = [];

				if (this.isError(res)) {
					return [];
				}

				let apps = res.result.Applications || [];
				for (let app of apps) {
					applications.push({
						"Hostname": app.Hostname,
						"Port": app.Port,
						"Username": app.Username,
						"Password": app.Password,
						"RemoteApp": this.cleanAppName(app.RemoteApp),
						"ConnectionName": app.ConnectionName
					});
				}

				return applications;
			});
	}

	getApplicationForUser(): angular.IPromise<IApplication[]> {
		return this.rpc.call({method: "ServiceApplications.GetListForCurrentUser", id: 1})
			.then((res: IRpcResponse): IApplication[] => {
				let applications: IApplication[] = [];

				if (this.isError(res)) {
					return [];
				}

				let apps = res.result.Applications || [];
				for (let app of apps) {
					applications.push({
						"Hostname": app.Hostname,
						"Port": app.Port,
						"Username": app.Username,
						"Password": app.Password,
						"RemoteApp": this.cleanAppName(app.RemoteApp),
						"ConnectionName": app.ConnectionName
					});
				}

				return applications;
			});
	}

	unpublish(application: IApplication): angular.IPromise<void> {
		return this.rpc.call({
			method: "ServiceApplications.UnpublishApplication",
			params: [{"ApplicationName": application.RemoteApp}],
			id: 1
		}).then((res: IRpcResponse): void => {
			this.isError(res);
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

	private cleanAppName(appName: string): string {
		if (appName) {
			return appName.replace(/^\|\|/, "");
		} else {
			return "Desktop";
		}
	}
}

angular.module("haptic.applications").service("ApplicationsSvc", ApplicationsSvc);
