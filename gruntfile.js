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

module.exports = function(grunt) {

	require("load-grunt-tasks")(grunt);

	grunt.initConfig({

		concat: {
			css: {
				src: [
					"bower_components/angular-material/angular-material.min.css",
					"bower_components/angular-ui-grid/ui-grid.min.css",
					"bower_components/angular-material-icons/angular-material-icons.css"
				],
				dest: "website/css/libs.min.css"
			}
		},

		copy: {
			lib: {
				src: [
					"bower_components/requirejs/require.js",
					"bower_components/jquery/dist/jquery.min.js", "bower_components/jquery/dist/jquery.min.map",
					"bower_components/lodash/lodash.min.js",
					"bower_components/angular/angular.min.js", "bower_components/angular/angular.min.js.map",
					"bower_components/angular-cookies/angular-cookies.min.js", "bower_components/angular-cookies/angular-cookies.min.js.map",
					"bower_components/angular-animate/angular-animate.min.js", "bower_components/angular-animate/angular-animate.min.js.map",
					"bower_components/angular-aria/angular-aria.min.js", "bower_components/angular-aria/angular-aria.min.js.map",
					"bower_components/angular-material/angular-material.min.js",
					"bower_components/angular-material-icons/angular-material-icons.min.js",
					"bower_components/angular-ui-grid/ui-grid.min.js",
					"bower_components/angular-ui-router/release/angular-ui-router.min.js",
					"bower_components/ui-router-extras/release/ct-ui-router-extras.min.js",
					"bower_components/ng-flow/dist/ng-flow-standalone.min.js"
				],
				dest: "website/js/lib/",
				expand: true,
				flatten: true
			},
			uigrid: {
				cwd: "bower_components/angular-ui-grid/",
				src: ["ui-grid.svg", "ui-grid.ttf", "ui-grid.eot", "ui-grid.woff"],
				dest: "website/css/",
				expand: true
			}
		},

		sync: {
			assets: {
				files: [{
					cwd: "website/ts/",
					src: ["**/*.html", "components/**/*.json"],
					dest: "website/js/"
				}],
				verbose: true
			}
		},

		less: {
			dist: {
				options: {
					compress: true
				},
				files: {
					"website/css/app.min.css": ["website/less/**/*.less"]
				}
			}
		},

		ts: {
			dist: {
				files: { "website/js/": ["website/ts/**/*.ts"] },
				tsconfig: "website/ts/tsconfig.json"
			},
			debug: {
				files: { "website/js/": ["website/ts/**/*.ts"] },
				tsconfig: "website/ts/tsconfig.json",
				watch: "website/ts/"
			}
		},

		tslint: {
			options: {
				configuration: grunt.file.readJSON("tslint.json")
			},
			files: {
				src: ["website/ts/**/*.ts"]
			}
		},

		watch: {
			ts: {
				files: ["website/ts/**/*.ts"],
				tasks: ["ts:debug"]
			},
			assets: {
				files: ["website/ts/**/*.html", "website/ts/components/**/*.json"],
				tasks: ["sync"]
			},
			less: {
				files: ["website/less/**/*.less"],
				tasks: ["less"]
			}
		}

	});

	grunt.registerTask("build-libs", ["concat", "copy:lib", "copy:uigrid"]);
	grunt.registerTask("build", ["tslint", "ts:dist", "less", "sync"]);
	grunt.registerTask("start", ["build", 	"watch"]);

};
