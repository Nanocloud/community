/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export interface IApplication {
		Id: number;
		Alias: string;
		CollectionName: string;
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
							Alias: "eclipse",
							CollectionName: "winapps",
							DisplayName: "Eclipse",
							IconContents: "",
							FilePath: "C:\\Program File\\Eclipse\\eclipse"
						}];
					}

					let apps = JSON.parse(res.result.ApplicationsJsonArray);
					console.log(apps);
					for (let app of apps) {
						applications.push({
							"Id": 1, // s TODO Have an ID
							"Alias": app.Alias,
							"CollectionName": app.CollectionName,
							"DisplayName": app.DisplayName,
							"IconContents": app.IconContents,
							"FilePath": app.FilePath
						});
					}
					console.log(applications);

					return applications;
				});
		}

		unpublish(application: IApplication): angular.IPromise<void> {
			return this.rpc.call({
				method: "ServiceApplications.UnpublishApplication",
				params: [{"ApplicationName": application.Alias}],
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
					.content(res.error.code === 0 ? "Internal Error" : res.error)
					.position("top right")
			);
			return true;
		}

	}

	app.service("ApplicationsService", ApplicationsService);
}
