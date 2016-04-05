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

// jshint mocha:true

var chai = require('chai');
var expect = chai.expect;
var sync = require('urllib-sync');
var extend = require('extend-object');
var plugins = require('./assertions/plugins');
var JSONAPIValidator = require('jsonapi-validator').Validator;
var clone = require('clone');

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';
extend(exports, plugins);

var expect = function(value) {
  if (plugins.chai === null) {
    exports.initialize();
  }
  if (value !== undefined && value !== null && value.then !== undefined) {
    var test = plugins.chai.expect(value).eventually;
    recordedExpects.push(test);
    return test;
  } else {
    return plugins.chai.expect(value);
  }
};

var nano = {
  PROTOCOL: process.env.NANOCLOUD_PROTOCOL || 'https',
  HOST: process.env.NANOCLOUD_HOST || 'localhost',
  PORT: process.env.NANOCLOUD_PORT || '443',
  CLIENTID: process.env.NANOCLOUD_CLIENTID || '9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae:9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341',
  ADMIN_USERNAME: process.env.NANOCLOUD_ADMIN_USERNAME || 'admin@nanocloud.com',
  ADMIN_PASSWORD: process.env.NANOCLOUD_ADMIN_PASSWORD || 'admin',
  USER_FIRSTNAME: process.env.NANOCLOUD_USER_FIRSTNAME || 'Nano',
  USER_LASTNAME: process.env.NANOCLOUD_USER_LASTNAME || 'Test',
  USER_EMAIL: process.env.NANOCLOUD_USER_EMAIL || 'test@nanocloud.com',
  USER_PASSWORD: process.env.NANOCLOUD_USER_PASSWORD || 'Nano@123',
  _request: function(user) {
    var makeRequest = function(verb, url, data, options) {
      var headers = (options && options.headers) ? options.headers : {};
      if (user) {
        headers.Authorization = 'Bearer ' + user.access_token;
        headers['Content-Type'] = 'application/json';
      }

      var request = sync.request(nano.PROTOCOL + '://' + nano.HOST + ':' + nano.PORT + '/' + url, {
        method: verb,
        data: data,
        headers: headers
      });

      // Get pure javascript object out of response Buffer
      var contentType = request.headers['content-type'].split(';');
      if (contentType.indexOf('application/json') !== -1 || contentType.indexOf('application/vnd.api+json') !== -1) {
        if (typeof request.data === 'string') {
          request.data = JSON.parse(request.data);
        } else if (typeof request.data === 'object') {
          request.data = JSON.parse(request.data.toString());
          if (typeof request.data !== 'object') {
          request.data = JSON.parse(request.data);
          }
        }
      }

      return request;
    };

    return {
      response: null,
      post: function(url, params, options) {
        this.response = makeRequest('POST', url, params, options);

        return this;
      },
      get: function(url, data) {
        this.response = makeRequest('GET', url, data || null);

        return this;
      },
      delete: function(url) {
        this. response = makeRequest('DELETE', url);

        return this;
      },
      shouldReturn: function(code) {
        it('should return ' + code, function() {
          expect(this).to.have.status(code);
        }.bind(this));

        return this;
      },
      shouldBeJSON: function() {
        it('should be valid JSON', function() {
          expect(this.response.headers).to.have.property('content-type');

          var values = this.response.headers['content-type'].split(';');
          expect(values).to.satisfy(function(type) {
           return (type.indexOf('application/json') !== -1 || type.indexOf('application/vnd.api+json') !== -1);
          });
        }.bind(this));

        return this;
      },
      shouldBeJSONAPI: function() {
        this.shouldBeJSON();
        var validator = new JSONAPIValidator();
        var valid = true;

        try {
          validator.validate(this.response.data);
        } catch (e) {
            valid = false;
        }


        it('should comply to JSON API schema', function() {
          expect(valid).to.equal(true);
        });

        return this;
      },
      shouldBeJSONAPIError: function(error) {

        this.shouldBeJSONAPI();

        it('should return expected json-api error schema', function() {
          var assert = new chai.Assertion(this.response.data.errors);

          assert.to.containSubset([error]);
        }.bind(this));

        return this;
      },
      shouldComplyToNotJsonAPI: function(schema) {

        it('should return expected schema', function() {
          expect(this).to.have.schema(schema);
        }.bind(this));

        return this;
      },
      shouldComplyTo: function(schema) {
        var JSONAPIschema = require('./JSONAPIschema.json');
        var clonedJSONAPIschema = clone(JSONAPIschema);
        clonedJSONAPIschema.definitions.attributes = schema;

        return this.shouldComplyToNotJsonAPI(clonedJSONAPIschema);
      }
    };
  },
  as: function(user) {
    return this._request(user);
  },
  get: function(url) {
    return this._request().get(url);
  },
  post: function(url, params, options) {
    return this._request().post(url, params, options);
  },
  login: function(credentials) {
    var user = sync.request(nano.PROTOCOL + '://' + nano.HOST + ':' + nano.PORT + '/oauth/token', {
      method: 'POST',
      dataType: 'json',
      data: {
        username: credentials.username,
        password: credentials.password,
        grant_type: 'password'
      },
      headers: {
        Authorization: 'Basic ' + new Buffer(this.CLIENTID).toString('base64')
      }
    });

    return user.data;
  }
};

module.exports = nano;
module.exports.expect = expect;
