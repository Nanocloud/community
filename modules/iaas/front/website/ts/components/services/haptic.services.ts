/// <reference path="../../../../typings/tsd.d.ts" />
/// <reference path="../../core.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
/// <amd-dependency path="angular-ui-grid" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "AmdTools";
import { MainMenu } from "MainMenu";

let componentName = "services";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future", "ui.grid", "ui.grid.expandable"]);

let states: angular.ui.IState[] = [{
	name: "admin.services",
	url: "/services",
	controller: "ServicesCtrl",
	controllerAs: "servicesCtrl",
	templateUrl: getTemplateUrl(componentName, "services.html")
}];

MainMenu.add({
	stateName: states[0].name,
	title: "Services",
	ico: "cloud"
});

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
