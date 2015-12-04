/// <reference path="../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-material" />
/// <amd-dependency path="angular-material-icons" />
import * as angular from "angular";

// create the main module
let app = angular.module("haptic", ["ngMaterial", "ngMdIcons"]);

app.config(["$mdThemingProvider", function($mdThemingProvider: angular.material.IThemingProvider) {

	$mdThemingProvider.theme("default").primaryPalette("blue");

}]);

let plugins: string[] = ["core", "login", "services", "users", "applications", "history", "presenter"]; // should be loaded via the backend

// load the available plugins to the main module
let deps: string[] = [];
for (var pn of plugins) {
	deps.push("components/" + pn + "/haptic." + pn);
	app.requires.push("haptic." + pn);
}
requirejs(deps, function() {
	// manually start up angular application
	angular.bootstrap(document, ["haptic"], { strictDi: true });
});
