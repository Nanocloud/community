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

    next(null, server);
  });
}

function getServers(project, ids, next) {

  var servers = [];
  var nova = new ostack.Nova(URL + ':8774/v2/' + PROJECT_ID, project.token);

  async.forEachOf(ids, function(id, key, callback) {
    nova.getServer(id, function(error, server) {
      if (error) {
        return next(error);
      }

      servers.push(server);
      callback();
    });
  }, function (error) {

    if (error) {
      return next(error);
    }

    return next(null, servers);
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

function associateFloatingIP(project, server, next) {

  var nova = new ostack.Nova(URL + ':8774/v2/' + PROJECT_ID, project.token);

  nova.listFloatingIps(function(error, floatingIPs) {

    if (error) {
      return next(error);
    }

    async.filter(floatingIPs, function(floatingIP, callback) {

      callback(!floatingIP.instance_id);
    }, function(availableIPs) {

      var _associateFloatingIP = function(server, ip, callback) {

        nova.associateFloatingIp(server.id, ip.ip, function(error) {

          if (error) {
            return callback(error);
          }

          callback(null, server);
        });

      };

      if (availableIPs.length == 0) {

        nova.createFloatingIp({}, function(error, result) {

          if (error) {
            next(error);
          }

          return _associateFloatingIP(server, result, next);
        });
      } else {
        var selectedIP = availableIPs[0];
        _associateFloatingIP(server, selectedIP, next);
      }

    });
  });
}

var project = null;

var provisionLinux = function(callback) {
  var linuxServer = null;

  async.waterfall([
    function(next) {

      createServer(project, {
        "name": "Bamboo Linux",
        "imageRef": URL + ":9292/v2/images/" + '7d771989-2ccb-47fb-bbb4-75ee6bd00f2f',
        "flavorRef": URL + ":8774/v2/flavors/2"
      }, function(error, _server) {

        if (error) {
          next(error);
        }

        linuxServer = _server;
        next(null);
      });
    },
    function waitForLinuxToBeOnline(next) {

      getServer(project, linuxServer.id, function(error, _server) {

        if (error) {
          return next(error);
        }

        if (_server.status != "ACTIVE") {
          return setTimeout(function() {
            waitForLinuxToBeOnline(next);
          }, 1000);
        }

        return next(null);
      });
    },
    function(next) {

      associateFloatingIP(project, linuxServer, function(error) {

        if (error) {
          next(error);
        }

        next(null);
      });
    }
  ], function(error) {

    if (error) {
      return callback(error);
    }

    console.log('Linux is online');
    return callback(null);
  });
};

var provisionWindows = function(callback) {
  var windowsServer = null;
  var image = null;

  async.waterfall([
    function(next) {
      uploadImage(project, next);
    },
    function(_image, next) {
      image = _image;

      createServer(project, {
        "name": "Bamboo Windows",
        "imageRef": URL + ":9292/v2/images/" + image.id,
        "flavorRef": URL + ":8774/v2/flavors/3"
      }, function(error, _server) {

        if (error) {
          return next(error);
        }

        windowsServer = _server;
        next(null);
      });
    },
    function waitForWindowsToBeOnline(next) {

      getServer(project, windowsServer.id, function(error, _server) {

        if (error) {
          return next(error);
        }

        if (_server.status != "ACTIVE") {
          return setTimeout(function() {
            waitForWindowsToBeOnline(next);
          }, 1000);
        }

        return next(null);
      });
    }
  ], function(error) {

    if (error) {
      return callback(error);
    }

    console.log('Windows is online');
    return callback(null);
  });
};

async.waterfall([
  login,
  getProject,
  function(_project, next) {
    project = _project;

    async.parallel([
      provisionLinux,
      provisionWindows
    ], function(error) {
      if (error) {
        next(error);
      }

      next(null, project);
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
