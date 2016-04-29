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

// Required environment variables
var LINUX_IMAGE_ID = process.env.DEPLOYMENT_LINUX_IMAGE_ID || null;
var PASSWORD = process.env.DEPLOYMENT_OS_PASSWORD || "";
var PROJECT_ID = process.env.DEPLOYMENT_OS_PROJECT_ID || '';
var USERNAME = process.env.DEPLOYMENT_OS_USERNAME || "";

// Optional environment variables
var INSTALLATION_SCRIPT = process.env.DEPLOYMENT_INSTALLATION_SCRIPT || './installCommunity.sh';
var INSTALL_SCRIPT_PATH = process.env.DEPLOYMENT_OS_INSTALL_SCRIPT_PATH || './installDocker.sh';
var KEY_NAME = process.env.DEPLOYMENT_OS_KEY_NAME || 'Bamboo';
var KEY_PATH = process.env.DEPLOYMENT_OS_KEY_PATH || './id_rsa';
var NETWORK_NAME= process.env.DEPLOYMENT_NETWORK_NAME || 'nano-net';
var PUBLIC_IP = process.env.DEPLOYMENT_PUBLIC_IP || null;
var SSH_PORT = process.env.DEPLOYMENT_OS_SSH_PORT || 22;
var LINUX_SECURITY_GROUPS = process.env.DEPLOYMENT_LINUX_SECURITY_GROUPS.split(';') || [
  "HTTP and HTTPS",
  "SSH"
];
var WINDOWS_IMAGE_ID = process.env.DEPLOYMENT_WINDOWS_IMAGE_ID || null;
var WINDOWS_IMAGE_PATH = process.env.DEPLOYMENT_OS_WINDOWS_IMAGE_PATH || './windows.qcow2';
var WINDOWS_SECURITY_GROUPS = process.env.DEPLOYMENT_WINDOWS_SECURITY_GROUPS.split(';') || [
  "Plaza",
  "LDAPS",
  "SSH",
  "RDP"
];
var LINUX_VM_NAME = process.env.DEPLOYMENT_LINUX_VM_NAME || 'Bamboo Linux';
var WINDOWS_VM_NAME = process.env.DEPLOYMENT_WINDOWS_VM_NAME || 'Bamboo Windows';

var nanoOS = require('./libnanoOpenstack');
var async = require('async');
var Promise = require('promise');
var test_port = require('test-port');

var URL = process.env.DEPLOYMENT_OS_URL || "http://openstack.nanocloud.org";

var project = null;
var _resolveWindowsIP = null;
var windowsIP = new Promise(function(resolve) {
  _resolveWindowsIP = resolve;
});

var provisionLinux = function(callback) {
  var linuxServer = null;
  var linuxIP = null;

  async.waterfall([
    function(next) { // Boot Linux server

      project.createServer({
        "name": LINUX_VM_NAME,
        "imageRef": URL + ":9292/v2/images/" + LINUX_IMAGE_ID,
        "flavorRef": URL + ":8774/v2/flavors/2",
        "key_name": KEY_NAME
      }, function(error, _server) {

        if (error) {
          next(error);
        }

        linuxServer = _server;
        next(null);
      });
    },
    function waitForLinuxToBeOnline(next) { // Wait for Linux to boot

      linuxServer.getStatus(function(error, status) {

        if (error) {
          return next(error);
        }

        if (status != "ACTIVE") {
          return setTimeout(function() {
            waitForLinuxToBeOnline(next);
          }, 1000);
        }

        return next(null);
      });
    },
    function(next) { // Open SSH and HTTPS port

      linuxServer.assignSecurityGroup(LINUX_SECURITY_GROUPS, function(error) {
        next(error);
      });
    },
    function(next) { // Associate public IP

      linuxServer.associateFloatingIP(PUBLIC_IP, function(error, ip) {

        if (error) {
          return next(error);
        }

        linuxIP = ip;
        next(null);
      });
    },
    function waitForSSHToBeAvailable(next) {

      test_port(SSH_PORT, linuxIP.ip, function(isOpen) {
        if (isOpen) {
          return next(null);
        }

        return setTimeout(function() {
          waitForSSHToBeAvailable(next);
        }, 1000);
      });
    },
    function(next) {

      linuxServer.execute(linuxIP, INSTALL_SCRIPT_PATH, KEY_PATH, function(error, response) {

        if (error) {
          return next(error);
        }

        next(null);
      });
    },
    function(next) {

      windowsIP.then(function(ip) {
        console.log('Windows IP + ' + ip);
        next(null, ip);
      });
    },
    function(winIP, next) { // Save Windows IP to be read back in the VM

      linuxServer.execute(linuxIP, "echo " + winIP + " > windowsIP", KEY_PATH, function(error) {
        next(error);
      });

    },
    function(next) { // Install community

      linuxServer.execute(linuxIP, INSTALLATION_SCRIPT, KEY_PATH, function(error, response) {

        if (error) {
          return next(error);
        }

        next(null, response);
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
    function(next) { // Upload qcow2
      if (WINDOWS_IMAGE_ID === null) {
        project.uploadImage(WINDOWS_IMAGE_PATH, {
          name: WINDOWS_VM_NAME + " image",
          visibility: 'private',
          disk_format: 'qcow2',
          container_format: 'bare'
        }, function(error, _image) {
          next(error, _image);
        });
      } else {
        next(null, null);
      }
    },
    function(_image, next) { // Boot Windows
      if (_image === null) {
        imageId = WINDOWS_IMAGE_ID;
      } else {
        imageId = _image.id;
      }

      project.createServer({
        "name": WINDOWS_VM_NAME,
        "imageRef": URL + ":9292/v2/images/" + imageId,
        "flavorRef": URL + ":8774/v2/flavors/3"
      }, function(error, _server) {

        if (error) {
          return next(error);
        }

        windowsServer = _server;
        next(null);
      });
    },
    function waitForWindowsToBeOnline(next) { // Wait for Windows to boot

      windowsServer.get(function(error, _server) {

        if (error) {
          return next(error);
        }

        if (_server.status != "ACTIVE") {
          return setTimeout(function() {
            waitForWindowsToBeOnline(next);
          }, 1000);
        }

        _resolveWindowsIP(_server.addresses[NETWORK_NAME][0].addr);
        return next(null);
      });
    },
    function (next) { // Assign right security group to Windows

      windowsServer.assignSecurityGroup(WINDOWS_SECURITY_GROUPS, function(error) {
        next(error);
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
  function(next) { // Login
    nanoOS.login(USERNAME, PASSWORD, function(error, user) {
      next(error, user);
    });
  },
  function (user, next) { // Get project
    user.getProject(PROJECT_ID, function(error, project) {
      next(error, project);
    });
  },
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
