/// <reference path='../../typings/tsd.d.ts' />

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
				title: "Home",
				url: "/",
				ico: "home"
			}, {
				title: "Users",
				url: "/users",
				ico: "people"
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
