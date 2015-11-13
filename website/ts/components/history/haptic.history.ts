/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
/// <amd-dependency path="angular-ui-grid" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "../core/services/AmdTools";
import { MainMenu } from "../core/services/MainMenu";

let componentName = "history";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future", "ui.grid", "ui.grid.expandable"]);

let states: angular.ui.IState[] = [{
	name: "admin.history",
	url: "/history",
	controller: "HistoryCtrl",
	controllerAs: "historyCtrl",
	templateUrl: getTemplateUrl(componentName, "stats.html")
}];

MainMenu.add({
	stateName: states[0].name,
	title: "History",
	ico: "equalizer"
});

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
