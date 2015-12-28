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

/// <reference path="../../../../typings/tsd.d.ts" />
/// <reference path="../../core.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
/// <amd-dependency path="angular-cookies" />
/// <amd-dependency path="angular-ui-grid" />
/// <amd-dependency path="ng-flow" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "AmdTools";
import { MainMenu } from "MainMenu";

let componentName = "applications";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future", "ui.grid", "ngCookies", "flow"]);

let states: angular.ui.IState[] = [{
	name: "admin.applications",
	url: "/applications",
	controller: "ApplicationsCtrl",
	controllerAs: "applicationsCtrl",
	templateUrl: getTemplateUrl(componentName, "applications.html")
}];

MainMenu.add({
	stateName: states[0].name,
	title: "Applications",
	ico: "apps"
});

app.config(["$controllerProvider", "$provide", "$futureStateProvider", "flowFactoryProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any,
	flowFactoryProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	registerCtrlFutureStates(componentName, $futureStateProvider, states);

	flowFactoryProvider.defaults = {
		headers: {
			"Authorization": "Bearer " + localStorage["accessToken"]
		}
	}

}]);
