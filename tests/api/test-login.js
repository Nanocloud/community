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
        type: {'type': 'string'}
      },
      additionalProperties: false
    };

    var request = nano.post('oauth/token', {
      username: 'admin@nanocloud.com',
      password: 'admin',
      grant_type: 'password'
    }, {
      headers: {
        Authorization: 'Basic ' + new Buffer(nano.CLIENTID).toString('base64'),
        'Content-Type': 'application/json'
      }
    }).shouldReturn(200)
        .shouldBeJSON()
        .shouldComplyToNotJsonAPI(expectedSchema);

    it('should issue Bearer tokens', function() {
      expect(request.response.data.type).to.equal('Bearer')
    })
  })
}
