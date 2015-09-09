/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class UserApplicationCtrl {

		applications: any;

		static $inject = [
			"ApplicationsService"
		];

		constructor(
			private applicationsSrv: ApplicationsService
		) {
			this.loadApplications();
		}

		loadApplications(): angular.IPromise<void> {
			return this.applicationsSrv.getApplicationForUser().then((applications: IApplication[]) => {
				this.applications = applications;
			});
		}

		openApplication(application: IApplication, e: MouseEvent) {
			window.open("/guacamole/client.xhtml?id=c%2F" + application.ConnectionName, "_blank");
		}
	}

	app.controller("UserApplicationCtrl", UserApplicationCtrl);
}
