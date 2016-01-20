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

		sync: {
			assets: {
				files: [{
					cwd: "website/ts/",
					src: ["**/*.html"],
					dest: "website/js/"
				}],
				verbose: true,
				compareUsing: "md5"
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
			assets: {
				files: ["website/ts/**/*.html"],
				tasks: ["sync"]
			},
			less: {
				files: ["website/less/**/*.less"],
				tasks: ["less"]
			}
		}

	});

	grunt.registerTask("build", ["tslint", "ts:dist", "sync"]);
	grunt.registerTask("watch-xy", ["sync", "watch"]);
	grunt.registerTask("watch-ts", ["tslint", "ts:debug"]);

};
