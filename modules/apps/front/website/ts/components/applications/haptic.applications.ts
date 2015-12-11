/// <reference path="../../../../typings/tsd.d.ts" />
/// <reference path="../../core.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
/// <amd-dependency path="angular-cookies" />
/// <amd-dependency path="angular-ui-grid" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "AmdTools";
import { MainMenu } from "MainMenu";

let componentName = "applications";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future", "ui.grid", "ngCookies"]);

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

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
