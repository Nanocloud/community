/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class LoginCtrl {

		credentials: any;
		isLoggedIn: boolean;

		static $inject = [
			"$location",
			"AuthenticationService"
		];

		constructor(
			private $location: angular.ILocationService,
			private authSrv: AuthenticationService
		) {
			this.credentials = {
				"email": "",
				"password": ""
			};
		}

		signIn(e: MouseEvent) {
			this.authSrv.authenticate(this.credentials).then((success: boolean) => {
				this.isLoggedIn = success;
			});
			if (this.isLoggedIn === true) {
				this.$location.path("/");
				window.location.href = "/";
			} else {
				this.$location.path("/login.html");
				window.location.href = "/login.html";
			}
		}
	}

	app.controller("LoginCtrl", LoginCtrl);
}
