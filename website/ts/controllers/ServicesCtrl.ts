/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class ServicesCtrl {

		services: any;
		colors: any;

		static $inject = [
			"ServicesService",
			"$mdDialog",
			"$interval"
		];

		constructor(
			private servicesSrv: ServicesService,
			private $mdDialog: angular.material.IDialogService,
			$interval: ng.IIntervalService
		) {
			this.colors = {
				downloading: "#4183D7",
				available: "#A2DED0",
				booting: "#EB9532",
				running: "#26A65B"
			};

			this.loadServices();
			$interval(
				this.loadServices.bind(this),
				5000
				);
		}

		loadServices(): angular.IPromise<void> {
			return this.servicesSrv.getAll().then((services: IService[]) => {
				this.services = services;
			});
		}

		startWindowsDownload(e: MouseEvent, service: IService) {
			let o = this.getDefaultServiceDlgOpt(e);
			o.locals = { service: service };
			return this.$mdDialog.show(o);
		}

		toggle(service: IService) {
			if (! service.Locked) {
				return this.servicesSrv.start(service);
			}
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
