/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export interface IApplication {
		Hostname: string;
		Port: string;
		Username: string;
		Password: string;
		RemoteApp: string;
		ConnectionName: string;
	}

	export class ApplicationsService {

		static $inject = [
			"RpcService",
			"$mdToast"
		];
		constructor(
			private rpc: RpcService,
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

					let apps = res.result.Applications;
					for (let app of apps) {
						applications.push({
							"Hostname": app.Hostname,
							"Port": app.Port,
							"Username": app.Username,
							"Password": app.Password,
							"RemoteApp": app.RemoteApp || "Desktop",
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
							"RemoteApp": app.RemoteApp || "Desktop",
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

	}

	app.service("ApplicationsService", ApplicationsService);
}
