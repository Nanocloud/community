/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class LoginCtrl {

		credentials: any;

		static $inject = [
			"$location",
			"AuthenticationService"
		];

		constructor(
			private $location: angular.ILocationService,
			private authSrv: AuthenticationService
		) {
			this.credentials = {
				"mail": "",
				"password": ""
			};
		}

		signIn(e: MouseEvent) {
			// let isLoggedIn = this.authSrv.authenticate(this.credentials);
			let isLoggedIn = this.authSrv.tmpAuth(this.credentials);
			if (isLoggedIn === true) {
				this.$location.path("/");
				window.location.href = "/";
			} else {
				this.$location.path("/login");
				window.location.href = "/login";
			}
		}
	}

	app.controller("LoginCtrl", LoginCtrl);
}
