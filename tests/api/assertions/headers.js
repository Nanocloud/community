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

module.exports = function (chai, utils) {

  utils.addMethod(chai.Assertion.prototype, 'header', function (key, expected) {

    var headerValue = this._obj.response.headers[key.toLowerCase()];

    if(arguments.length === 1) {
      this.assert(
        headerValue !== undefined && headerValue !== null,
        'expected header '+ key +' to exist',
        'expected header '+ key +' not to exist'
      );
    } else if (expected instanceof RegExp) {
      this.assert(
        expected.test(headerValue),
        'expected header '+ key + ' with value ' + headerValue + ' to match regex '+expected,
        'expected header '+ key + ' with value ' + headerValue + ' not to match regex '+expected
      );
    } else if (typeof(expected) === 'function') {
      expected(headerValue);
    } else {
      this.assert(
        headerValue === expected,
        'expected header '+ key + ' with value ' + headerValue + ' to match '+expected,
        'expected header '+ key + ' with value ' + headerValue + ' not to match '+expected
      );
    }
  });
};
