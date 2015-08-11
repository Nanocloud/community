/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class ServicesCtrl {

		gridOptions: any;
		services: any;

		static $inject = [
			"ServicesService",
		];

		constructor(
			private servicesSrv: ServicesService,
			private $mdDialog: angular.material.IDialogService
		) {
			this.gridOptions = {
				data: [],
				rowHeight: 36,
				columnDefs: [
					{ field: "Ico" },
					{ field: "Name" },
				]
			};

			this.loadServices();
		}

		loadServices(): angular.IPromise<void> {
			return this.servicesSrv.getAll().then((services: IService[]) => {
				this.services = services;
			});
		}

	}

	app.controller("ServicesCtrl", ServicesCtrl);
}
