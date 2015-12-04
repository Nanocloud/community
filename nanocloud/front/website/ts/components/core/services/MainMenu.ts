
export interface INavMenu {
	stateName?: string;
	title?: string;
	ico?: string;
}

export module MainMenu {

	export var menus: INavMenu[] = [];

	export function add(menu: INavMenu) {
		menus.push(menu);
	}

}