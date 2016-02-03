/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

/// <reference path="../../../../../typings/tsd.d.ts" />
/// <amd-dependency path="../services/HistorySvc" />
import { HistorySvc, IHistoryInfo } from "../services/HistorySvc";

"use strict";

// missing prop in the tsd file
interface IExpandableRowScope {
	subGridVariable: string;
}
interface IGridOptions extends uiGrid.IGridOptionsOf<IHistoryInfo> {
	expandableRowScope: IExpandableRowScope;
}

export class HistoryCtrl {

	gridOptions: IGridOptions;

	static $inject = [
		"HistorySvc",
		"$mdDialog"
	];
	constructor(
		private historySvc: HistorySvc,
		private $mdDialog: angular.material.IDialogService
	) {
		this.gridOptions = {
			expandableRowTemplate: "./js/components/history/views/stat.html",
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

	setData(hist: IHistoryInfo[]) {
		let data: any[] = [];
		for (let stat of hist) {
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
			data.push(s);
		}

		this.gridOptions.data = data;
	}

	loadStats(): angular.IPromise<void> {
		return this.historySvc.getAll().then((hist: IHistoryInfo[]) => {
			this.setData(hist);
		});
	}
	
}

angular.module("haptic.history").controller("HistoryCtrl", HistoryCtrl);
