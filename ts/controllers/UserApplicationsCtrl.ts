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
	}

	app.controller("UserApplicationCtrl", UserApplicationCtrl);
}
