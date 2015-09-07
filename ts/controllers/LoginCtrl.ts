/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class LoginCtrl {

		credentials: any;

		static $inject = [
			"$location",
			"AuthenticationService",
			"$mdToast"
		];

		constructor(
			private $location: angular.ILocationService,
			private authSrv: AuthenticationService,
			private $mdToast: angular.material.IToastService
		) {
			this.credentials = {
				"email": "",
				"password": ""
			};
		}

		signIn(e: MouseEvent) {
			this.authSrv.authenticate(this.credentials).then(
					(success: boolean) => {
						this.$location.path("/");
						window.location.href = "/";
					},
					(error: boolean) => {
						this.$mdToast.show(
								this.$mdToast.simple()
								.content("Authentication failed: Email or Password incorrect")
								.position("top right")
								);
					}
			);
		}
	}

	app.controller("LoginCtrl", LoginCtrl);
}
