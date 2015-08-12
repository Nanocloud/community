/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export class AuthenticationService {

		static $inject = [
			"RpcService",
		];
		constructor(
			private rpc: RpcService,
			private $mdToast: angular.material.IToastService
		) {
		}

		authenticate(credentials): angular.IPromise<boolean> {
			return this.rpc.call({
				method: "Auth.Authenticate",
				id: 1,
				mail: credentials.mail,
				password: credentials.password
			}).then((res: IRpcResponse): boolean => {
					if (this.isError(res)) {
						return true;
					}
					return true;
			});
		}

		tmpAuth(credentials): boolean {
			return true;
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

	app.service("AuthenticationService", AuthenticationService);
}
