(function() {
  "use strict";
  window.addEventListener("load", init);
  // let time1;
  // let time2;
  function init() {
    // createDevice(null, null);
    // deviceLogin();
    // deviceInfo();
    deviceDisconnect();
  }

  function deviceDisconnect() {
    createDevice(connectDevice, disconnect);
  }
  function deviceInfo() {
    createDevice(getInfo, null);
  }

  function disconnect(data, bearer) {
    let time1 = performance.now()
    let settings = {
      "async": true,
      "crossDomain": true,
      "url": "https://api.bfranzen.me/disconnect",
      "method": "DELETE",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": bearer
        // "cache-control": "no-cache",
        // "Postman-Token": "be187d81-5a0f-4958-b507-de958572128d"
      },
      "processData": false,
      "data": data
    }
    $.ajax(settings).done(function (response) {
      let time2 = performance.now();
      let bd = document.querySelector("body");
      let p = document.createElement("p");
      p.setAttribute("id", "disconnect-time");
      p.innerText = "Time to disconnect is " + (time2 - time1) / 1000 + " seconds";
      bd.appendChild(p);
      console.log(response);
    });
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
      let bd = document.querySelector("body");
      let p = document.createElement("p");
      p.setAttribute("id", "getinfo-time");
      p.innerText = "Getting device info took " + ((time2 - time1)/1000) + " seconds.";
      bd.appendChild(p);
    });
  }
  // Creates a device then logs in
  function deviceLogin() {
    createDevice(connectDevice);
  }

  // Makes ajax call to the connect endpoint to connect/login to a device
  function connectDevice(bearer, deviceName, nextCall) {
    let time1 = performance.now();
    let data = "{\n\t\"name\": \""+ deviceName + "\",\n\t\"latitude\": 123,\n\t\"longitude\": 234,\n\t\"email\": \"tesst@test\",\n\t\"phone\": \"12345\",\n\t\"password\": \"1234567\",\n\t\"passwordConf\": \"1234567\"\n}";
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
      "data": data
    }

    $.ajax(settings).done(function (response) {
      let time2 = performance.now();
      console.log(response);
      if (nextCall === null) {
        let bd = document.querySelector("body");
        let p = document.createElement("p");
        p.setAttribute("id", "connect-time");
        p.innerText = "Logging into a device took " + ((time2 - time1)/1000) + " seconds.";
        bd.appendChild(p);

      } else {
        nextCall(data, bearer);
      }
    });
  }

  // Sends ajax request to create a device, if login is true
  // it will use credentials to login
  function createDevice(func, nextCall) {
    let baseName = "timeDisconnect1";
    let time1 = performance.now();
    // let time2 = performance.now();
    let numDevices = 1;
    for (let i = 0; i < numDevices; i++) {
      let deviceName = "timeDisconnect5";
      let settings = {
        "async": true,
        "crossDomain": true,
        "url": "https://api.bfranzen.me/device",
        "method": "POST",
        "headers": {
          "Content-Type": "application/json",
        },
        "processData": false,
        "data": "{\n\t\"name\": \""+ deviceName +"\",\n\t\"latitude\": 123,\n\t\"longitude\": 234,\n\t\"email\": \"tesst@test\",\n\t\"phone\": \"12345\",\n\t\"password\": \"1234567\",\n\t\"passwordConf\": \"1234567\"\n}"
      }
      let jqxhr = $.ajax(settings)
      jqxhr.done(function (response) {
        let time2 = performance.now();
        console.log(response);
        let bearer = jqxhr.getResponseHeader("Authorization");
        if (func !== null) {
          func(bearer, deviceName, nextCall);
        } else {
          let bd = document.querySelector("body");
          let p = document.createElement("p");
          p.setAttribute("id", "createdevice-time");
          p.innerText = "Creating a device took " + ((time2 - time1)/1000) + " seconds.";
          bd.appendChild(p);
        }
      });
    }
  }
})();
