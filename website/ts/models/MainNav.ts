/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend.models {
	
	export interface INavMenu {
		title: string;
		url: string;
		ico: string;
	}
	
	export class MainNav {
		
		constructor() {
			this.current = this.menus[0];
		}
		
		menus: INavMenu[] = [
			{
				title: "Services",
				url: "/",
				ico: "cloud"
			}, {
				title: "Applications",
				url: "/applications",
				ico: "apps"
			}, {
				title: "Users",
				url: "/users",
				ico: "people"
				/*
			}, {
				title: "Stats",
				url: "/stats",
				ico: "equalizer"
			}, {
				title: "Capacity Planning",
				url: "/capacity_planning",
				ico: "trending_up"
				*/
			}
		];
		
		private _current: INavMenu;
		get current(): INavMenu {
			return this._current;
		}
		set current(menu: INavMenu) {
			this._current = menu;
		}
		
	}
	
}
