/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
/// <amd-dependency path="angular-cookies" />
/// <amd-dependency path="angular-ui-grid" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "../core/services/AmdTools";

let componentName = "applications";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future", "ui.grid", "ngCookies"]);

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	let states: angular.ui.IState[] = [{
		name: "admin.applications",
		url: "applications",
		controller: "ApplicationsCtrl",
		controllerAs: "applicationsCtrl",
		templateUrl: getTemplateUrl(componentName, "applications.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
