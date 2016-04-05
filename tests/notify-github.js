#!/usr/bin/nodejs
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

var proc = require('child_process');

/*
 * Example params:
 * {
 *   github_password,
 *   pull_sha,
 *   state,
 *   description,
 *   context,
 *   url
 * }
 */
var notify = function(params) {

  var command = 'curl -H "Authorization: token ' + params.github_password + '" --request POST -k --data \'{\"state\": \"'+ params.state +'\", \"description\": \"'+ params.description +'\", \"context\": \"'+ params.context +'\", \"target_url\": \"' + params.url + '\"}\' https://api.github.com/repos/nanocloud/community/statuses/' + params.pull_sha;

    console.log("Notify github with command: " + command)
    proc.exec(command, function (err, stdout, stderr) {
	if (err) {
	    console.log("Error contacting gitHub: " + stderr);
	} else {
	    console.log("Github notify completed")
	    console.log(stdout)
	}
    });
}

var _usage = false;

var usage = function() {
  if (_usage == false) {
    console.log("Usage: notify-github: pull_sha state description context url");
    console.log("Expects github_password to be set in environment");

    _usage = true;
  }

  return _usage;
}

// If directly invoked from command line
// Usage: notify-github: pull_sha state description context url
// Expects github_password to be set in environment
if (require.main === module) {

  var pull_sha = process.argv[2] || usage.call();
  var state = process.argv[3] || usage.call();
  var description = process.argv[4] || usage.call();
  var context = process.argv[5] || usage.call();
  var url = process.argv[6] || usage.call();
  var github_password = global.process.env["github_password"] || usage();

  if (_usage == false) {
    notify({
      github_password: github_password,
      pull_sha: pull_sha,
      state: state,
      description: description,
      context: context,
      url: url
    });
  }
}

module.exports.notify = notify;
