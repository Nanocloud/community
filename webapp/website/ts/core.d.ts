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
