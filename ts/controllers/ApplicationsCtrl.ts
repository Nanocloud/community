/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class ApplicationsCtrl {

		gridOptions: any;

		static $inject = [
			"ApplicationsService",
		];

		constructor(
			private applicationsSrv: ApplicationsService,
			private $mdDialog: angular.material.IDialogService
		) {
			this.gridOptions = {
				data: [],
				rowHeight: 36,
				columnDefs: [
					{
						name: "Icon",
						field: "IconContents",
						cellTemplate: "<img src=\"data:image/jpeg;base64,{{grid.getCellValue(row, col)}}\">"
					},
					{ field: "DisplayName" },
					{ field: "FilePath" },
					{
						name: "edit",
						displayName: "",
						enableColumnMenu: false,
						cellTemplate: "\
							<md-button ng-click='grid.appScope.applicationsCtrl.startEditApplication($event, row.entity)'>\
								<ng-md-icon icon='edit' size='14'></ng-md-icon> Edit\
							</md-button>\
							<md-button ng-click='grid.appScope.applicationsCtrl.startDeleteApplication($event, row.entity)'>\
								<ng-md-icon icon='delete' size='14'></ng-md-icon> Delete\
							</md-button>"
					}
				]	
			};

			this.loadApplications();
		}

		get applications(): IApplication[] {
			return this.gridOptions.data;
		}
		set applications(value: IApplication[]) {
			this.gridOptions.data = value;
		}

		loadApplications(): angular.IPromise<void> {
			return this.applicationsSrv.getAll().then((applications: IApplication[]) => {
				this.applications = applications;
			});
		}

	}

	app.controller("ApplicationsCtrl", ApplicationsCtrl);
}
