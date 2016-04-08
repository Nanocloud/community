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

  downloadToken = null;
  pathToBeTested = "C:\\";

  describe("Get download token", function() {

    var requestToken = nano.as(admin).get("api/files/token", { filename : pathToBeTested + "\\Windows\\system.ini" })
      .shouldReturn(200)

    it("should return a token",  function() {
      return expect(requestToken.response.data.token).not.to.be.empty; 
    });

    if (requestToken.response.data.token) {
      downloadToken = requestToken.response.data.token;
    }
  });

  describe("Download file", function() {

    var downloadFile = nano.as(admin).get("api/files", { filename : pathToBeTested + "\\Windows\\system.ini" , token : downloadToken  })
      .shouldReturn(200)

    it("should download the file",  function() {
      return expect(downloadFile.response.data).not.to.be.empty; 
    });

    it("should match system.ini content",  function() {
      return expect(downloadFile.response.data.toString()).to.have.string('drivers');
    });
  });
}
