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

  describe("List sessions with no sessions", function() {


    var request = nano.as(admin).get("api/sessions")
    .shouldReturn(200)
    .shouldBeJSONAPI();

    it("shoud be an empty array", function() {
      return expect(request.response.data.data).to.be.empty;
    });

  });
};
