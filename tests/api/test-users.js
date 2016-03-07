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

  describe("List users", function() {

    var expectedSchema = {
      type: 'object',
      properties: {
        email: {type: 'string'},
        activated: {type: 'boolean'},
        is_admin: {type: 'boolean'},
        first_name: {type: 'string'},
        last_name: {type: 'string'},
        sam: {type: 'string'},
        windows_password: {type: 'string'},
      },
      required: ['email', 'first_name', 'activated', 'is_admin', 'first_name', 'last_name', 'sam', 'windows_password'],
      additionalProperties: false
    };

    var request = nano.as(admin).get("api/users")
        .shouldReturn(200)
        .shouldBeJSONAPI()
        .shouldComplyTo(expectedSchema);

    it("should contain the admin",  function() {
      return expect(request).to.comprise.of.json({
        email: 'admin@nanocloud.com',
        activated: true,
        is_admin: true,
        first_name: "John",
        last_name: "Doe",
        sam: "Administrator",
        windows_password: "Nanocloud123+"
      });
    })
  })

  describe("Create user", function() {

    var request = nano.as(admin).post("api/users", {
      "data" : {
        "type": "user",
        "attributes": {
          "first_name": "{{TEST_USER_FIRSTNAME}}",
          "last_name": "{{TEST_USER_LASTNAME}}",
          "email": "{{TEST_USER_EMAIL}}",
          "password": "{{TEST_USER_PASSWORD}}"
        }
      }
    }).shouldReturn(200)
        .shouldBeJSON();

    it("Should return 403", function() {
      expect(request).to.have.status(403);
    })
  })
}
