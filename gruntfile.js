module.exports = function(grunt) {

	require("load-grunt-tasks")(grunt);

	grunt.initConfig({

		concat: {
			js: {
				src: [
					"bower_components/jquery/dist/jquery.min.js",
					"bower_components/lodash/lodash.min.js",
					"bower_components/angular/angular.min.js",
					"bower_components/angular-route/angular-route.min.js",
					"bower_components/angular-animate/angular-animate.min.js",
					"bower_components/angular-aria/angular-aria.min.js",
					"bower_components/angular-material/angular-material.min.js",
					"bower_components/angular-material-icons/angular-material-icons.min.js",
					"bower_components/angular-ui-grid/ui-grid.min.js",
					"bower_components/ng-flow/dist/ng-flow-standalone.js"
				],
				dest: "website/js/libs.min.js"
			},
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
			dist: {
				cwd: "bower_components/angular-ui-grid/",
				src: ["ui-grid.svg", "ui-grid.ttf", "ui-grid.eot", "ui-grid.woff"],
				dest: "website/css/",
				expand: true
			}	
		},

		less: {
			dist: {
				options: {
					compress: true
				},
				files: {
					"website/css/app.min.css": ["less/**/*.less"]
				}
			}
		},

		ts: {
			dist: {
				files: {
					"website/js/app.min.js": ["ts/**/*.ts"]
				},
				options: {
					target: "es5",
					fast: "never"
				}
			}
		},

		tslint: {
			options: {
				configuration: grunt.file.readJSON("tslint.json")
			},
			files: {
				src: ["ts/**/*.ts"]
			}
		},

		watch: {
			less: {
				files: [ "less/**/*.less" ],
				tasks: [ "less" ]
			},
			ts: {
				files: [ "ts/**/*.ts" ],
				tasks: [ "ts" ]
			}
		}

	});

	grunt.registerTask("update-libs", ["concat", "copy"]);
	grunt.registerTask("build", ["less", "tslint", "ts"]);

};