/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class ServicesCtrl {

		services: any;

		static $inject = [
			"ServicesService",
			"$mdDialog"
		];

		constructor(
			private servicesSrv: ServicesService,
			private $mdDialog: angular.material.IDialogService
		) {
			this.loadServices();
		}

		loadServices(): angular.IPromise<void> {
			return this.servicesSrv.getAll().then((services: IService[]) => {
				this.services = services;
			});
		}

		displayInfo(e: MouseEvent, service: IService) {
			let o = this.getDefaultServiceDlgOpt(e);
			o.locals = { service: service };
			return this.$mdDialog.show(o);
		}

		private getDefaultServiceDlgOpt(e: MouseEvent): angular.material.IDialogOptions {
			return {
				controller: "ServiceCtrl",
				controllerAs: "serviceCtrl",
				templateUrl: "./views/service.html",
				parent: angular.element(document.body),
				targetEvent: e
			};
		}


	}

	app.controller("ServicesCtrl", ServicesCtrl);
}
