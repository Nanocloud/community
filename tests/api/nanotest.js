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

var expect = require('chai').expect;
var sync = require('urllib-sync');
var extend = require('extend-object');
var plugins = require('./assertions/plugins');
var JSONAPIValidator = require('jsonapi-validator').Validator;

process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";
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
  _request: function(user) {
    var makeRequest = function(verb, url, data, options) {
      var headers = (options && options.headers) ? options.headers : {};
      if (user) {
        headers['Authorization'] = 'Bearer ' + user.access_token;
        headers['Content-Type'] = 'application/json';
      }

      var request = sync.request('https://localhost/' + url, {
        method: verb,
        data: data,
        headers: headers
      });

      // Get pure javascript object out of response Buffer
      if (request.headers['content-type'].split(';').indexOf('application/json') != -1) {
        request.data = JSON.parse(request.data.toString());
      }

      return request;
    };

    return {
      response: null,
      post: function(url, params, options) {
        this.response = makeRequest('POST', url, params, options);

        return this;
      },
      get: function(url) {
        this.response = makeRequest('GET', url);

        return this;
      },
      shouldReturn: function(code) {
        it("should return " + code, function() {
          expect(this).to.have.status(code)
        }.bind(this));

        return this;
      },
      shouldBeJSON: function() {
        it("should be valid JSON", function() {
          expect(this.response.headers).to.have.property('content-type');

          var values = this.response.headers['content-type'].split(';');
          expect(values).to.include('application/json')
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

        if (valid && !validator.isValid(this.response.data)) {
          valid = false;
        }

        it("should comply to JSON API schema", function() {
          expect(valid).to.equal(true);
        });

        return this;
      },
      shouldComplyToNotJsonAPI: function(schema) {

        it("should return expected schema", function() {
          expect(this).to.have.schema(schema)
        }.bind(this));

        return this;
      },
      shouldComplyTo: function(schema) {
        var JSONAPIschema = require('./JSONAPIschema.json');

        JSONAPIschema.definitions.attributes = schema;

        return this.shouldComplyToNotJsonAPI(JSONAPIschema);
      }
    }
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
    var clientId = '9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae:9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341';
    var user = sync.request('https://localhost/oauth/token', {
      method: 'POST',
      dataType: 'json',
      data: {
        username: credentials.username,
        password: credentials.password,
        grant_type: "password"
      },
      headers: {
        Authorization: 'Basic ' + new Buffer(clientId).toString('base64'),
        'Content-Type': 'application/json'
      }
    })

    return user.data;
  }
}

module.exports = nano
module.exports.expect = expect;
