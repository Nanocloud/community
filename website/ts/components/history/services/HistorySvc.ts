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

export interface IHistoryInfo {
	UserId: string;
	ConnectionId: string;
	Stats: IHistoryAtom[];
}

export interface IHistoryAtom {
	StartDate: string;
	EndDate: string;
}

export class HistorySvc {

	static $inject = [
		"RpcSvc",
		"$mdToast"
	];
	constructor(
		private rpc: RpcSvc,
		private $mdToast: angular.material.IToastService
	) {

	}

	getAll(): angular.IPromise<IHistoryInfo[]> {
		return this.rpc.call({ method: "ServiceHistory.GetList", id: 1 })
			.then((res: IRpcResponse) => {
				if (this.isError(res) || res.result.Histories === null) {
					return [];
				}
				return res.result.Histories;
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

angular.module("haptic.history").service("HistorySvc", HistorySvc);
