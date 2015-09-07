/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export class AuthenticationService {

		static $inject = [
			"$http"
			"$mdToast"
		];
		constructor(
			private $http: angular.IHttpService,
			private $mdToast: angular.material.IToastService
		) {

			alert(this.$http);

		}

		authenticate(credentials): angular.IPromise<boolean> {
			return this.$http.post("/login", {
				email: credentials.email,
				password: credentials.password
			})
			.then(function (response): boolean {
				return true;
			}, function (response): boolean {
				return false;
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

	app.service("AuthenticationService", AuthenticationService);
}
