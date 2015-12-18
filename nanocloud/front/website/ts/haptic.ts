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
