(function() {
  "use strict";
  window.addEventListener("load", init);
  // let time1;
  // let time2;
  function init() {
    // createDevice(null);
    // deviceLogin();
    deviceInfo();

  }

  function deviceInfo() {
    createDevice(getInfo);
  }

  function getInfo(bearer, deviceName) {
    let time1 = performance.now();
    var settings = {
      "async": true,
      "crossDomain": true,
      "url": "https://api.bfranzen.me/device-info/",
      "method": "GET",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": bearer
      },
      "processData": false,
      "data": "{\n\t\"name\": \""+deviceName+"\",\n\t\"password\": \"1234567\",\n\t\"passwordConf\": \"1234567\"\n}"
    }

    $.ajax(settings).done(function (response) {
      let time2 = performance.now();
      console.log(response);
      let p = document.getElementById("time-info");
      p.innerHTML = "Getting device info took " + ((time2 - time1)/1000) + " seconds.";
    });
  }
  // Creates a device then logs in
  function deviceLogin() {
    createDevice(connectDevice);
  }

  // Makes ajax call to the connect endpoint to connect/login to a device
  function connectDevice(bearer, deviceName) {
    let time1 = performance.now();
    let settings = {
      "async": true,
      "crossDomain": true,
      "url": "https://api.bfranzen.me/connect",
      "method": "POST",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": bearer
        // "cache-control": "no-cache",
        // "Postman-Token": "a666feca-010b-4179-8bf4-a23c5b4846bd"
      },
      "processData": false,
      "data": "{\n\t\"name\": \""+ deviceName + "\",\n\t\"latitude\": 123,\n\t\"longitude\": 234,\n\t\"email\": \"tesst@test\",\n\t\"phone\": \"12345\",\n\t\"password\": \"1234567\",\n\t\"passwordConf\": \"1234567\"\n}"    }

    $.ajax(settings).done(function (response) {
      let time2 = performance.now();
      console.log(response);
      let p = document.getElementById("time-info");
      p.innerHTML = "Logging into a device took " + ((time2 - time1)/1000) + " seconds.";
    });
  }

  // Sends ajax request to create a device, if login is true
  // it will use credentials to login
  function createDevice(func) {
    let baseName = "loginTest";
    let time1 = performance.now();
    // let time2 = performance.now();
    let numDevices = 1;
    for (let i = 0; i < numDevices; i++) {
      let deviceName = "testDeviceInfo05";
      let settings = {
        "async": true,
        "crossDomain": true,
        "url": "https://api.bfranzen.me/device",
        "method": "POST",
        "headers": {
          "Content-Type": "application/json",
        },
        "processData": false,
        "data": "{\n\t\"name\": \""+deviceName+"\",\n\t\"latitude\": 123,\n\t\"longitude\": 234,\n\t\"email\": \"tesst@test\",\n\t\"phone\": \"12345\",\n\t\"password\": \"1234567\",\n\t\"passwordConf\": \"1234567\"\n}"
      }
      let jqxhr = $.ajax(settings)
      jqxhr.done(function (response) {
        let time2 = performance.now();
        console.log(response);
        let bearer = jqxhr.getResponseHeader("Authorization");
        if (func !== null) {
          func(bearer, deviceName);
        } else {
          let p = document.getElementById("time-info");
          p.innerHTML = "Creating "+ numDevices + " devices took " + ((time2 - time1)/1000) + " seconds.";
        }
      });
    }
  }
})();
