
declare module "AmdTools" {
	function overrideModuleRegisterer(app: angular.IModule, $controllerProvider: angular.IControllerProvider, $provide: angular.auto.IProvideService): void;
	function registerCtrlFutureStates(comptName: string, $futureStateProvider: any, states: angular.ui.IState[]): void;
	function getTemplateUrl(comptName: string, url: string): string;
}

declare module "MainMenu" {
	export interface INavMenu {
		stateName?: string;
		title?: string;
		ico?: string;
	}
	export module MainMenu {
		let menus: INavMenu[];
		function add(menu: any): void;
	}
}
