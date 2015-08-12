/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export interface IService {
		Id: number;
		Name: string;
		Ico: string;
	}

	export class ServicesService {

		static $inject = [
			"RpcService",
			"$mdToast"
		];
		constructor(
			private rpc: RpcService,
			private $mdToast: angular.material.IToastService
		) {

		}

		getAll(): angular.IPromise<IService[]> {
			return this.rpc.call({ method: "ServiceServices.GetList", id: 1 })
				.then((res: IRpcResponse): IService[] => {
					let services: IService[] = [];

					if (this.isError(res)) {
						return [{ // fake data
							Id: 1,
							Name: "Proxy",
							Ico: "public"
						}];
					}

					for (let srv of res.result.ServicesJsonArray) {
						services.push(JSON.parse(srv));
					}
					return services;
				});
		}

		private isError(res: IRpcResponse): boolean {
			if (res.error == null) {
				return false;
			}
			this.$mdToast.show(
				this.$mdToast.simple()
					.content(res.error.code === 0 ? "Internal Error" : res.error.message)
					.position("top right")
			);
			return true;
		}

	}

	app.service("ServicesService", ServicesService);
}
