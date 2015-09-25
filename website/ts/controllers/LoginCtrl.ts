/// <reference path='../../../typings/tsd.d.ts' />

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
			let user = {
				"Email": this.credentials.email
			};
			sessionStorage.setItem("user", user.Email);

			this.authSrv.authenticate(this.credentials).then(
					(response: any) => {
						if (typeof response.data === "string") {
							if (response.data.indexOf instanceof Function &&
									response.data.indexOf("<body layout=\"row\" ng-controller=\"MainCtrl as mainCtrl\">") !== -1) {
								this.$location.path("/admin.html");
								window.location.href = "/admin.html";
								return;
							}
						}
						this.$location.path("/");
						window.location.href = "/";
					},
					(error: any) => {
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
