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

var nano = require('./nanotest');
var expect = nano.expect;

module.exports = function(admin) {

  var expectedSchema = {
    type: 'object',
    properties: {
      email: {type: 'string'},
      activated: {type: 'boolean'},
      'is-admin': {type: 'boolean'},
      'first-name': {type: 'string'},
      'last-name': {type: 'string'},
    },
    required: ['email', 'first-name', 'activated', 'is-admin', 'first-name', 'last-name'],
    additionalProperties: false
  };

  describe('List users', function() {

    var request = nano.as(admin).get('api/users')
        .shouldReturn(200)
        .shouldBeJSONAPI()
        .shouldComplyTo(expectedSchema);

    it('should contain the admin',  function() {
      return expect(request).to.comprise.of.json({
        email: 'admin@nanocloud.com',
        activated: true,
        'is-admin': true,
        'first-name': 'Admin',
        'last-name': 'Nanocloud'
      });
    });
  });

  var user_id = null;
  describe('Create user', function() {

    var request = nano.as(admin).post('api/users', {
      'data' : {
        'type': 'user',
        'attributes': {
          'first-name': nano.USER_FIRSTNAME,
          'last-name': nano.USER_LASTNAME,
          'email': nano.USER_EMAIL,
          'password': nano.USER_PASSWORD
        }
      }
    }).shouldReturn(201)
    .shouldBeJSONAPI()
    .shouldComplyTo(expectedSchema);

    if (request.response.data.data) {
      user_id = request.response.data.data.id;
    }
  });

  describe('Remove user', function() {

    nano.as(admin).delete('api/users/' + user_id)
    .shouldReturn(200)
    .shouldBeJSONAPI();

  });
};
