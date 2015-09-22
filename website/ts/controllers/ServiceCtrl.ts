/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";
	
	class ServiceCtrl {

		service: IService;

		static $inject = [
			"ServicesService",
			"$mdDialog",
			"service"
		];
		constructor(
			private servicesSrv: ServicesService,
			private $mdDialog: angular.material.IDialogService,
			service: IService
		) {
			if (service) {
				this.service = angular.copy(service);
			}
		}

		accept(): void {
			if (this.servicesSrv.downloadStarted === false) {
				this.servicesSrv.download();
				this.$mdDialog.hide(this.service);
			} else {
				this.$mdDialog.cancel();
			}
		}

		close(): void {
			this.$mdDialog.cancel();
		}

	}

	app.controller("ServiceCtrl", ServiceCtrl);
}
