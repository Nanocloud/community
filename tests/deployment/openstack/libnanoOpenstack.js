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

var ostack = require('openstack-wrapper');
var async = require('async');
var fs = require('fs');
var ssh = require('ssh-exec');

function _getKeystone() {

  if (_getKeystone._keystone === undefined) {
    _getKeystone._keystone = new ostack.Keystone(URL + ':5000/v3/');
  }

  return _getKeystone._keystone;
}

var nano = {

  NanoOSProject: function(project) {

    this.getProject = function() {
      return project;
    };
  },

  NanoOSServer: function(project, server) {

    this.ip = null;
    this.getProject = function() {
      return project;
    };
    this.getServer = function() {
      return server;
    };
  },

  NanoOSUser: function(OSUser) {

    this.getUser = function() {
      return OSUser;
    };
  },
  login: function(username, password, callback) {
    _getKeystone().getToken(username, password, function(error, token) {

      if (error) {
        return callback(error);
      }

      callback(null, new nano.NanoOSUser(token));
    });
  }
};

nano.NanoOSUser.prototype.getProject = function(projectID, callback) {

  _getKeystone().getProjectToken(this.getUser().token, projectID, function(error, projectToken) {

    if (error) {
      return callback(error);
    }

    return callback(null, new nano.NanoOSProject(projectToken));
  }.bind(this));
};

nano.NanoOSProject.prototype._getNova = function() {

  if (this._nova === undefined) {
    this._nova = new ostack.Nova(URL + ':8774/v2/' + this.getProject().project.id, this.getProject().token);
  }

  return this._nova;
};

nano.NanoOSProject.prototype.createServer = function(data, callback) {

  this._getNova().createServer({
    server: data
  }, function(error, server) {
    callback(error, new nano.NanoOSServer(this, server));
  }.bind(this));
};

nano.NanoOSProject.prototype.getServer = function(id, callback) {

  this._getNova().getServer(id, function(error, server) {

    if (error) {
      callback(error);
    }

    callback(null, new nano.NanoOSServer(this, server));
  });
};

nano.NanoOSProject.prototype.getIP = function(id, callback) {

  this._getNova().getFloatingIp(id, function(error, ip) {

    if (error) {
      callback(error);
    }

    callback(null, ip);
  });
};

nano.NanoOSProject.prototype._getGlance = function() {

  if (this._glance === undefined) {
    this._glance = new ostack.Glance(URL + ':9292/v2/', this.getProject().token);
  }

  return this._glance;
};

nano.NanoOSProject.prototype.uploadImage = function(path, metadata, callback) {

  this._getGlance().queueImage(metadata, function(error, image) {

    if (error) {
      return callback(error);
    }

    var file = fs.createReadStream(path);
    this._getGlance().uploadImage(image.id, file, function(error) {

      if (error) {
        return callback(error);
      }

      return callback(null, image);
    });
  }.bind(this));

};

nano.NanoOSServer.prototype.get = function(callback) {

  this.getProject()._getNova().getServer(this.getServer().id, function(error, server) {

    if (error) {
      callback(error);
    }

    callback(null, server);
  });
};

nano.NanoOSServer.prototype.getStatus = function(callback) {

  this.get(function(error, server) {

    if (error) {
      callback(error);
    }

    callback(null, server.status);
  });
};

nano.NanoOSServer.prototype.assignSecurityGroup = function(groupName, callback) {

  if (typeof groupName === 'string') {
    groupName = [groupName];
  }

  async.each(groupName, function(group, next) {
    this.getProject()._getNova().assignSecurityGroup(group, this.getServer().id, function(error) {
      next(error);
    });
  }.bind(this), function(err) {
    callback(err);
  });
};

nano.NanoOSServer.prototype.execute = function(ip, scriptPath, keyPath, callback) {

  var user = 'debian';

  fs.readFile(scriptPath, function(err, script) {

    if (err) {
      callback(err);
    }

    ssh(script, {
      user: user,
      host: ip.ip,
      key: keyPath
    }).pipe(process.stdout);
  });
};

nano.NanoOSServer.prototype.associateFloatingIP = function(callback) {

  this.getProject()._getNova().listFloatingIps(function(error, floatingIPs) {

    if (error) {
      return callback(error);
    }

    async.filter(floatingIPs, function(floatingIP, callback) {
      callback(!floatingIP.instance_id);
    }, function(availableIPs) {

      var _associateFloatingIP = function(ip, callback) {

        this.getProject()._getNova().associateFloatingIp(this.getServer().id, ip.ip, function(error) {

          if (error) {
            return callback(error);
          }

          callback(null, ip);
        });
      }.bind(this);

      if (availableIPs.length == 0) {

        this.getProject()._getNova().createFloatingIp({}, function(error, result) {

          if (error) {
            return callback(error);
          }

          return _associateFloatingIP(result, callback);
        });
      } else {
        var selectedIP = availableIPs[0];
        _associateFloatingIP(selectedIP, callback);
      }
    }.bind(this));
  }.bind(this));
};

module.exports = nano;
