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

var tv4 = require('tv4');

module.exports = function (chai, utils) {

  utils.addMethod(chai.Assertion.prototype, 'schema', function (schema) {

    var object = this._obj.response.data;

    var valid = tv4.validate(object, schema);

    var composeErrorMessage = function () {
      var errorMsg = 'expected body to match JSON schema ' + JSON.stringify(schema) + '.';
      if(tv4.error !== null) {
        errorMsg += '\n error: ' + tv4.error.message + '.\n data path: ' + tv4.error.dataPath + '.\n schema path: ' + tv4.error.schemaPath + '.';
      }
      return errorMsg;
    };

    this.assert(
      valid,
      composeErrorMessage(),
      'expected body to not match JSON schema ' + JSON.stringify(schema)
    );
  });
};
