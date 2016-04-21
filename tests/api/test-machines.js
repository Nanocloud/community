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

module.exports = function(admin) {

  var expectedSchema = {
    type: 'object',
    properties: {
      name: {type: 'string'},
      ip: {type: 'string'},
      type: {type: 'string'},
      status: {type: 'string'},
      'admin-password': {type: 'string'},
      platform: {type: 'string'},
      progress: {type: 'string'},
    },
    required: ['name', 'ip', 'type', 'status', 'platform', 'progress'],
    additionalProperties: false
  };

  describe('List machines', function() {
    nano.as(admin).get('api/machines')
        .shouldReturn(200)
        .shouldBeJSONAPI()
        .shouldComplyTo(expectedSchema);
  });
};
