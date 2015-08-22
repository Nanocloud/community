/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";
	
	class ServiceCtrl {
		
		service: IService;
		
		static $inject = [
			"$mdDialog",
			"service"
		];
		constructor(
			private $mdDialog: angular.material.IDialogService,
			service: IService
		) {
			if (service) {
				this.service = angular.copy(service);
			}
		}

		close(): void {
			this.$mdDialog.cancel();
		}

		start(): void {
			this.$mdDialog.cancel();
		}

		stop(): void {
			this.$mdDialog.cancel();
		}

	}

	app.controller("ServiceCtrl", ServiceCtrl);
}
