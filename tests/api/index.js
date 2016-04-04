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

describe("nanocloud is Online", function() {
  var request = nano.get('').shouldReturn(200);
});

var admin = nano.login({
  username: nano.ADMIN_USERNAME,
  password: nano.ADMIN_PASSWORD
});

describe("Windows should be up", function() {
  var request = nano.as(admin).get('api/machines')
      .shouldReturn(200)
      .shouldBeJSONAPI();

  it('Should have one Windows up', function() {
    expect(request.response.data.data).to.exist;
    expect(request.response.data.data).to.have.lengthOf(1);
    expect(request.response.data.data[0].attributes.status).to.equal('up');
  });
});

require('./test-login')();
require('./test-users')(admin);
require('./test-sessions')(admin);
