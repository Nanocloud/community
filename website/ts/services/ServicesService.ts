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

/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export interface IService {
		Ico: string;
		Name: string;
		DisplayName: string;
		Locked: boolean;
		Status: string;
		FontColor: string;
		VM: string;
	}

	export class ServicesService {

		static $inject = [
			"RpcService",
			"$mdToast"
		];

		downloadStarted = false;
		windowsReady = false;

		constructor(
			private rpc: RpcService,
			private $mdToast: angular.material.IToastService
		) {

		}

		getAll(): angular.IPromise<IService[]> {
			return this.rpc.call({ method: "ServiceIaas.GetList", id: 1 })
				.then((res: IRpcResponse): IService[] => {
					let services: IService[] = [];

					if (this.isError(res) || res.result.VmListJsonArray === undefined) {
						return [];
					} else {
						for (let srv of JSON.parse(res.result.VmListJsonArray)) {
							if (srv.Ico === "windows") {
								this.downloadStarted = false;
								this.windowsReady = true;
							} else {
								this.windowsReady = false;
							}
							services.push(srv);
						}
						return services;
					}
				});
		}

		download(): void  {
			this.rpc.call({ method: "ServiceIaas.Download", id: 1 })
				.then((res: IRpcResponse): void => {
					if (! this.isError(res) && res.result.Success === true) {
						this.downloadStarted = true;
					}
				});
		}

		start(service: IService): void  {
			this.rpc.call({ method: "ServiceIaas.Start", id: 1, params: [{"vmName": service.Name}] })
				.then((res: IRpcResponse): void => {
					if (! this.isError(res) && res.result.Success === true) {
						service.Status = "booting";
					} else {
						service.Status = "available";
					}
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

	app.service("ServicesService", ServicesService);
}
