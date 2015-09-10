/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";
	
	class ServiceCtrl {

		static $inject = [
			"ServicesService",
			"$mdDialog"
		];
		constructor(
			private servicesSrv: ServicesService,
			private $mdDialog: angular.material.IDialogService
		) {
		}

		accept(): void {
			if (this.servicesSrv.downloadStarted === false) {
				this.servicesSrv.download();
			}
			this.$mdDialog.cancel();
		}

		close(): void {
			this.$mdDialog.cancel();
		}

	}

	app.controller("ServiceCtrl", ServiceCtrl);
}
