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

var URL = process.env.DEPLOYMENT_OS_URL || "http://openstack.nanocloud.org";
var USERNAME = process.env.DEPLOYMENT_OS_USERNAME || "";
var PASSWORD = process.env.DEPLOYMENT_OS_PASSWORD || "";
var PROJECT_ID = process.env.DEPLOYMENT_OS_PROJECT_ID || '';

var ostack = require('openstack-wrapper');
var keystone = new ostack.Keystone(URL + ':5000/v3/');
var async = require('async');
var fs = require('fs');

function login(next) {
  keystone.getToken(USERNAME, PASSWORD, function(error, token) {

    if (error) {
      return next(error);
    }

    next(null, token);
  });
}

function getProject(user, next) {
  keystone.getProjectToken(user.token, PROJECT_ID, function(error, project_token) {

    if (error) {
      return next(error);
    }

    return next(null, project_token);
  });
}

function createServer(project, server, next) {

  var nova = new ostack.Nova(URL + ':8774/v2/' + PROJECT_ID, project.token);
  nova.createServer({
    "server": server
  }, function(error, server) {

    if (error) {
      next(error);
    }

    next(null, server);
  });
}

function listServers(project, next) {

  var nova = new ostack.Nova(URL + ':8774/v2/' + PROJECT_ID, project.token);
  nova.listServers(function(error, list) {
    if (error) {
      return next(error);
    }

    return next(null, list);
  });
}

function getServer(project, id, next) {

  var nova = new ostack.Nova(URL + ':8774/v2/' + PROJECT_ID, project.token);
  nova.getServer(id, function(error, server) {
    if (error) {
      return next(error);
    }

    return next(null, server);
  });
}

function uploadImage(project, next) {

  var glance = new ostack.Glance(URL + ':9292/v2/', project.token);

  glance.queueImage({
    name: "bamboo",
    visibility: 'private',
    disk_format: 'qcow2',
    container_format: 'bare'
  }, function(error, image) {

    if (error) {
      return next(error);
    }

    var file = fs.createReadStream('./windows.qcow2');
    glance.uploadImage(image.id, file, function(error) {

      if (error) {
        return next(error);
      }

      return next(null, image);
    });
  });
}

var project = null;
var image = null;
var server = null;

async.waterfall([
  login,
  getProject,
  function(_project, next) {
    project = _project;

    next(null, project);
  },
  uploadImage,
  function(_image, next) {
    image = _image;

    createServer(project, {
      "name": "Spawned by Bamboo",
      "imageRef": URL + ":9292/v2/images/" + image.id,
      "flavorRef": URL + ":8774/v2/flavors/3"
    }, function(error, _server) {

      if (error) {
        next(error);
      }

      server = _server;
      next(null, project);
    });
  },
  function waitForServerToBeOnline(project, next) { // Wait for server to start

    getServer(project, server.id, function(error, _server) {

      if (error) {
        next(error);
      }

      if (_server.status != "ACTIVE") {
        setTimeout(function() {
          waitForServerToBeOnline(project, next);
        }, 1000);
      } else {
        return next(null, _server);
      }
    });
  }
], function(error) {

  if (error) {
    console.log(error);
    console.log(error.stack);
    return ;
  }
  console.log('Deployement complete');
});
