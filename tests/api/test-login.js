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

var nano = require('./nanotest');
var expect = nano.expect;

module.exports = function() {

  describe('Login as an admin', function() {

    var expectedSchema = {
      type: 'object',
      properties: {
        access_token: {type: 'string'},
        token_type: {'type': 'string'},
				expires_in: {type: 'integer'}
      },
      required: ['access_token', 'token_type', 'expires_in'],
      additionalProperties: false
    };

    var request = nano.post('oauth/token', {
      username: nano.ADMIN_USERNAME,
      password: nano.ADMIN_PASSWORD,
      grant_type: 'password'
    }, {
      headers: {
        Authorization: 'Basic ' + new Buffer(nano.CLIENTID).toString('base64')
      }
    }).shouldReturn(200)
        .shouldBeJSON()
        .shouldComplyToNotJsonAPI(expectedSchema);

    it('should issue Bearer tokens', function() {
      expect(request.response.data.token_type).to.equal('Bearer')
    })
  });

  var expectedErrorSchema = {
    type: 'object',
    properties: {
      error: {type: 'string'},
      error_description: {'type': 'string'}
    },
    required: ['error', 'error_description'],
    additionalProperties: false
  };

  describe('Login with an invalid user', function() {

    var request = nano.post('oauth/token', {
      username: 'george.burdell@gatech.edu',
      password: 'fake',
      grant_type: 'password'
    }, {
      headers: {
        Authorization: 'Basic ' + new Buffer(nano.CLIENTID).toString('base64')
      }
    }).shouldReturn(401)
        .shouldBeJSON()
        .shouldComplyToNotJsonAPI(expectedErrorSchema);

    it('should return access_denied', function() {
      expect(request.response.data.error).to.equal("access_denied");
    });

    it('should return "Invalid User Credentials" as an error description', function() {
      expect(request.response.data.error_description).to.equal("Invalid User Credentials");
    });
  });

  describe('Login without specifying grant_type', function() {

    var request = nano.post('oauth/token', {
      username: 'george.burdell@gatech.edu',
      password: 'fake',
    }, {
      headers: {
        Authorization: 'Basic ' + new Buffer(nano.CLIENTID).toString('base64'),
      }
    }).shouldReturn(400)
        .shouldBeJSON()
        .shouldComplyToNotJsonAPI(expectedErrorSchema);

    it('should return access_denied', function() {
      expect(request.response.data.error).to.equal("invalid_request");
    });

    it('should return "Invalid User Credentials" as an error description', function() {
      expect(request.response.data.error_description).to.equal("grant_type is missing");
    });
  })
}
