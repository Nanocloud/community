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

/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend.models {
	
	export interface INavMenu {
		title: string;
		url: string;
		ico: string;
	}
	
	export class MainNav {

		HasPanelServices = false;
		HasPanelApplications 	 = false;
		HasPanelUsers 			 = true;
		HasPanelStats 			 = true;
		HasPanelCapacityPlanning = false;

		private setupMenu() {

			if (this.HasPanelServices) {
				this.menus.push({
					title: "Services",
					url: "/",
					ico: "cloud"
				});
			}

			if (this.HasPanelApplications) {
				this.menus.push({
					title: "Applications",
					url: "/applications",
					ico: "apps"
				});
			}

			if (this.HasPanelUsers) {
				this.menus.push({
					title: "Users",
					url: "/users",
					ico: "people"
				});
			}

			if (this.HasPanelStats) {
				this.menus.push({
					title: "Stats",
					url: "/stats",
					ico: "equalizer"
				});
			}

			if (this.HasPanelCapacityPlanning) {
				this.menus.push({
					title: "Capacity Planning",
					url: "/capacity_planning",
					ico: "trending_up"
				});
			}
		}		
		
		constructor() {

			this.setupMenu();

			this.current = this.menus[0];
		}
		
		menus: INavMenu[] = [];
		
		private _current: INavMenu;
		get current(): INavMenu {
			return this._current;
		}
		set current(menu: INavMenu) {
			this._current = menu;
		}

	}
	
}
