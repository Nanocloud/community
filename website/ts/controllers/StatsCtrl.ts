/// <reference path='../../../typings/tsd.d.ts' />

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
				data: [],
				rowHeight: 36,
				columnDefs: [
					{ field: "Firstname" },
					{ field: "Lastname" },
					{ field: "Email" }
				]	
			};
			
			this.loadStats();
		}

		get stats(): IStat[] {
			return this.gridOptions.data;
		}
		set stats(value: IStat[]) {
			this.gridOptions.data = value;
		}

		loadStats(): angular.IPromise<void> {
			return this.statsSrv.getAll().then((stats: IStat[]) => {
				this.stats = stats;
			});
		}
		
		addStat(stat: IStat): void {
			this.stats.push(stat);
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

	app.controller("StatsCtrl", StatsCtrl);
}
