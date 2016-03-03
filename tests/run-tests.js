#!/usr/bin/nodejs

var http = require('https');
var proc = require('child_process');
var async = require("async");
var HOST = process.argv[2];
var TOKEN = null;

process.env.NODE_TLS_REJECT_UNAUTHORIZED = "0";

function request(options, callback) {
  var message = "";
  var status = null;

  var headers = options.headers || {};
  var param = JSON.stringify(options.param) || "";

  if (TOKEN) {
    headers.Authorization = "Bearer " + TOKEN
  }

  var options = {
    host: HOST,
    path: options.path,
    method: options.verb,
    port: 443,
    headers: headers
  };

  var req = http.request(options, function(res) {
    status = res.statusCode;

    res.setEncoding('utf8');
    res.on('data', function(chunk) {
      message += chunk;
    });
    res.on('end', function() {
      console.log("Call to " + options.path + " returned " + status)
      if (callback && status == 200)
        callback(JSON.parse(message));
      if (status != 200)
        console.log(message);
    })
  });

  req.on('error', function(e) {
    console.log('problem with request: ' + e.message);
  });

  req.write(param);
  req.end();
}

function waitForNanocloudToBeOnline(next) {
  var command = 'curl --output /dev/null --insecure --silent --write-out \'%{http_code}\n\' "https://'+ HOST +'"';

  console.log("Try to connect")
  proc.exec(command, function (err, stdout, stderr) {

    console.log(stdout);
    if (!err) {
      if (stdout == "200\n") {
        console.log("Nanocloud available");

        return next();
      }
    }

    setTimeout(function () {
      waitForNanocloudToBeOnline(next);
    }, 2000);
  });

}

function bootWindows(next) {
  console.log("Booting Windows");
  request({
    path: '/api/iaas/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64/start',
    verb: 'POST'
  }, function() {
    if (next)
      next()
  });
}

function waitForWindowsToBeRunning(next) {
  console.log("Waiting for Windows to be running....");
  request({
    path: '/api/iaas',
    verb: 'GET',
  }, function(res) {
    if (res.data[0].attributes.status != "running") {
      waitForWindowsToBeRunning(next);
    }
    else
      next();
  });
}

function login(next) {
  request({
    path: '/oauth/token',
    verb: 'POST',
    param: {
      username: "admin@nanocloud.com",
      password: "admin",
      grant_type: "password"
    },
    headers: {
      'Content-Type': 'application/json',
      'Authorization': "Basic OTQwNWZiNmIwZTU5ZDI5OTdlM2M3NzdhMjJkOGYwZTYxN2E5ZjViMzZiNjU2NWM3NTc5ZTViZTZkZWI4ZjdhZTo5MDUwZDY3YzJiZTA5NDNmMmM2MzUwNzA1MmRkZWRiM2FlMzRhMzBlMzliYmJiZGFiMjQxYzkzZjhiNWNmMzQx"
    }
  }, function(res) {
    TOKEN = res.access_token;
    console.log("Got token : " + TOKEN);
    next();
  });
}

function setHostInEnv(next) {
  var command = 'sed -i "s/value\\": \\"127.0.0.1\\"/value\\": \\"'+ HOST +'\\"/g" api/NanoEnv.postman_environment'

  console.log("Setting host in api file")
  proc.exec(command, function (err, stdout, stderr) {

    if (err) {
      console.log(stdout);
      console.log(stderr);
      throw "Cannot set host in environment file"
    }
    return next();

  });
}

function done() {
  console.log('Ready to perform tests')
}

async.waterfall([
  waitForNanocloudToBeOnline,
  setHostInEnv,
  login,
  bootWindows,
  waitForWindowsToBeRunning,
  done
])
