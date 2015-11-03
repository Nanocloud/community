/// <reference path="../../../../../typings/tsd.d.ts" />
/// <amd-dependency path="../../core/services/RpcSvc" />
import { RpcSvc, IRpcResponse } from "../../core/services/RpcSvc";

"use strict";

export interface IStat {
	ConnectionId: string;
	StartDate: string;
	EndDate: string;
	Stats: any;
}

export class StatsSvc {

	static $inject = [
		"RpcSvc",
		"$mdToast"
	];
	constructor(
		private rpc: RpcSvc,
		private $mdToast: angular.material.IToastService
	) {

	}

	getAll(): angular.IPromise<IStat[]> {
		return this.rpc.call({ method: "ServiceHistory.GetList", id: 1 })
			.then((res: IRpcResponse): IStat[] => {
				let stats: IStat[] = [];

				if (this.isError(res)) {
					return [];
				}

				if (res.result.Histories) {
					for (let s of res.result.Histories) {
						stats.push(s);
					}
				}
				return stats;
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

angular.module("haptic.history").service("StatsSvc", StatsSvc);
