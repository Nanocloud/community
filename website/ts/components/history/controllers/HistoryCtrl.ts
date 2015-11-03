/// <reference path='../../../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";
	
	class StatsCtrl {

		gridOptions: any;

		static $inject = [
			"StatsService",
			"$mdDialog"
		];
		constructor(
			private statsSrv: StatsService,
			private $mdDialog: angular.material.IDialogService
		) {
			this.gridOptions = {
				expandableRowTemplate: "views/stat.html",
				expandableRowScope: {
					subGridVariable: "Stats"
				},
				columnDefs: [
					{
						displayName: "Connection Name",
						field: "ConnectionId"
					}
				]
			};

			this.loadStats();
		}

		get stats(): IStat[] {
			return this.gridOptions.data;
		}
		set stats(value: IStat[]) {
			let stats: any[] = [];
			for (var stat of value) {
				let s = {
					ConnectionId: stat.ConnectionId,
					subgridOptions: {
						columnDefs: [
							{ field: "StartDate" },
							{ field: "EndDate" }
						],
						data: stat.Stats
					}
				};
				stats.push(s);
			}

			this.gridOptions.data = stats;
		}

		loadStats(): angular.IPromise<void> {
			return this.statsSrv.getAll().then((stats: IStat[]) => {
				this.stats = stats;
			});
		}
		
		private getDefaultStatDlgOpt(e: MouseEvent): angular.material.IDialogOptions {
			return {
				controller: "StatsCtrl",
				controllerAs: "statsCtrl",
				templateUrl: "./views/stats.html",
				parent: angular.element(document.body),
				targetEvent: e
			};
		}
	}

	angular.module("haptic.history").controller("StatsCtrl", StatsCtrl);
}
