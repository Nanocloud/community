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

/// <reference path='../../typings/tsd.d.ts' />

requirejs.config({
	baseUrl: "/js/",
	paths: {
		"jquery": "lib/jquery.min",
		"lodash": "lib/lodash.min",
		"angular": "lib/angular.min",
		"angular-cookies": "lib/angular-cookies.min",
		"angular-animate": "lib/angular-animate.min",
		"angular-aria": "lib/angular-aria.min",
		"angular-material": "lib/angular-material.min",
		"angular-material-icons": "lib/angular-material-icons.min",
		"angular-ui-grid": "lib/ui-grid.min",
		"angular-ui-route": "lib/angular-ui-router.min",
		"angular-ui-router-extras": "lib/ct-ui-router-extras.min",
		"ng-flow": "lib/ng-flow-standalone.min",
		
		"AmdTools": "components/core/services/AmdTools",
		"MainMenu": "components/core/services/MainMenu"
	},
	shim: {
		"angular": { exports: "angular", deps: ["jquery"] },
		"angular-route": { deps: ["angular"] },
		"angular-aria": { deps: ["angular"] },
		"angular-animate": { deps: ["angular"] },
		"angular-cookies": { deps: ["angular"] },
		"angular-material": { deps: ["angular", "angular-animate", "angular-aria"] },
		"angular-material-icons": { deps: ["angular-material"] },
		"angular-ui-grid": { deps: ["angular"] },
		"angular-ui-route": { deps: ["angular"] },
		"angular-ui-router-extras": { deps: ["angular-ui-route"] },
		"ng-flow": { deps: ["angular"] }
	},
	deps: ["haptic"],
	waitSeconds: 25
});
