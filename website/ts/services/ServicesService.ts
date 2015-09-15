/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export interface IService {
		Ico: string;
		Name: string;
		DisplayName: string;
		Locked: boolean;
		Running: boolean;
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

		downloadStatus(): void  {
			this.rpc.call({ method: "ServiceIaas.DownloadStatus", id: 1 })
				.then((res: IRpcResponse): void => {
					if (! this.isError(res) && res.result.Success === true) {
						this.downloadStarted = true;
					} else {
						this.downloadStarted = false;
					}
				});
		}

		start(service: IService): void  {
			this.rpc.call({ method: "ServiceIaas.Start", id: 1, params: [{"vmName": service.Name}] })
				.then((res: IRpcResponse): void => {
					if (! this.isError(res) && res.result.Success === true) {
						service.Running = true;
					} else {
						service.Running = false;
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
