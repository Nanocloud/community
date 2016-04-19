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

module.exports = function(admin) {

  var fakeConnection = {
    'user-id': "fake-id",
    'connection-id': "fake-connection-id",
    'start-date': "fake-start-date",
    'end-date': "fake-end-date" 
  }

  var expectedSchema = {
    type: 'object',
    properties: {
      'user-id': {type: 'string'},
      'connection-id': {type: 'string'},
      'start-date': {type: 'string'},
      'end-date': {type: 'string'},
    },
    required: ['user-id', 'connection-id', 'start-date', 'end-date'],
    additionalProperties: false
  };

  describe("Create fake history entry", function() {

    var request = nano.as(admin).post('api/histories', {
      'data' : {
        'type': 'object',
        'attributes': {
          'user-id': fakeConnection["user-id"],
          'connection-id': fakeConnection["connection-id"],
          'start-date': fakeConnection["start-date"],
          'end-date': fakeConnection["end-date"]
        }
      }
    })
    .shouldReturn(201)
    .shouldBeJSONAPI()
    .shouldComplyTo(expectedSchema);
  });


  describe("List history entries", function() {

    var request = nano.as(admin).get('api/histories')
    .shouldReturn(200)
    .shouldBeJSONAPI()
    .shouldComplyTo(expectedSchema);

    it("should contains fake history entry",  function() {
      return expect(request).to.comprise.of.json({
        "user-id": "fake-id",
        "connection-id": "fake-connection-id",
        "start-date": "fake-start-date",
        "end-date": "fake-end-date",
      });
    });
  });
}
