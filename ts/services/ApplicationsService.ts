/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export interface IApplication {
		Id: number;
		DisplayName: string;
		IconContents: string;
		FilePath: string;
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
						// s TODO Display error
						return [{ // fake data
							Id: 1,
							DisplayName: "Eclipse",
							IconContents: "",
							FilePath: "C:\\Program File\\Eclipse\\eclipse"
						}];
					}

					let apps = JSON.parse(res.result.ApplicationsJsonArray);
					for (let app of apps) {
						applications.push({
							"Id": 1, // s TODO Have an ID
							"DisplayName": app.DisplayName,
							"IconContents": app.IconContents,
							"FilePath": app.FilePath
						});
					}
					return applications;
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

	app.service("ApplicationsService", ApplicationsService);
}
