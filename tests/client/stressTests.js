(function() {
  "use strict";

  const APP_SERVER_KEY = "BN6oGHmUe7MPtJNrpJzWSPjm-Iy3HmRo1TuvNcKgsGuwCBYYXDjXrM8r5wvFRdZO0kEnct_TDaX4sGTdIarLrJg";

//   Public Key:
// BN6oGHmUe7MPtJNrpJzWSPjm-Iy3HmRo1TuvNcKgsGuwCBYYXDjXrM8r5wvFRdZO0kEnct_TDaX4sGTdIarLrJg
//
// Private Key:
// um6cg6CWjqFo3Evs4xkejSca4BhxYfRBfiTRsCRGZy0
  window.addEventListener("load", init);
  // let time1;
  // let time2;
  function init() {
    // first parameter is name of device
    // second parameter is the call to be made after create deviceInfo
    // third one is for the function to call after the second function call
    // createDevice("t1", null, null);
    // deviceLogin("t2");
    // deviceInfo("t3");
    // deviceDisconnect("t4");
    // updateDevice("t5");
    document.getElementById("notify-me").addEventListener("click", notifyMe);
    document.getElementById("get-location").addEventListener("click", getLocation);

  }

  async function notifyMe() {
    // Register Service Worker
    console.log("Registering service worker...");
    const register = await navigator.serviceWorker.register("/service-worker.js", {
      scope: "/"
    }).then(console.error);
    console.log("Service Worker Registered...");

    // Register Push
    console.log("Registering Push...");
    const subscription = await register.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(APP_SERVER_KEY)
    }).then(console.error);
    console.log("Push Registered...");

    // Send Push Notification
    console.log("Sending Push...");
    await fetch("/subscribe", {
      method: "POST",
      body: JSON.stringify(subscription),
      headers: {
        "content-type": "application/json"
      }
    }).then(console.error);
    console.log("Push Sent...");
  }


  function urlBase64ToUint8Array(base64String) {
    const padding = "=".repeat((4 - base64String.length % 4) % 4);
    const base64 = (base64String + padding)
      .replace(/\-/g, "+")
      .replace(/_/g, "/");

    const rawData = window.atob(base64);
    const outputArray = new Uint8Array(rawData.length);

    for (let i = 0; i < rawData.length; ++i) {
      outputArray[i] = rawData.charCodeAt(i);
    }
    return outputArray;
  }


  function getLocation() {
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(function(pos) {
        position(pos);
      }, function(err) {
        // permissionDenied(err);
        // getLocation();
        console.error(err);
      }, {timeout:10000});
    }
  }

  function position(position) {
    let pos = {
      lat: position.coords.latitude,
      lng: position.coords.longitude
    };
    console.log(pos);
  }

  function deviceDisconnect(name) {
    createDevice(name, connectDevice, disconnect);
  }
  function deviceInfo(name) {
    createDevice(name, getInfo, null);
  }

  function updateDevice(name) {
    createDevice(name, updateDeviceInfo, null);
  }
  function updateDeviceInfo(bearer, deviceName, nextCall) {
    let time1 = performance.now();
    let settings = {
      "async": true,
      "crossDomain": true,
      "url": "https://api.bfranzen.me/device-info",
      "method": "PATCH",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": bearer,
        // "User-Agent": "PostmanRuntime/7.11.0",
        // "Accept": "*/*",
        // "Cache-Control": "no-cache",
        // "Postman-Token": "266e2fc8-94a4-4c8e-8e09-a37c89790339,0fd799d6-145d-40bc-a04c-6f9954db5be8",
        // "Host": "api.bfranzen.me",
        // "accept-encoding": "gzip, deflate",
        // "content-length": "170",
        // "Connection": "keep-alive",
        // "cache-control": "no-cache"
      },
      "processData": false,
      "data": "{\n\t\"name\": \"test_changeeeee\",\n\t\"latitude\": 333,\n\t\"longitude\": 333,\n\t\"email\": \"testingchange@changechange.com\",\n\t\"phone\": \"\",\n\t\"status\": \"\",\n\t\"oldPassword\": \"\",\n\t\"password\": \"\",\n\t\"passwordConf\": \"\"\n}"
    }

    $.ajax(settings).done(function (response) {
      let time2 = performance.now();
      let bd = document.querySelector("body");
      let p = document.createElement("p");
      p.setAttribute("id", "update-time");
      p.innerText = "Time to update device info is " + (time2 - time1) / 1000 + " seconds";
      bd.appendChild(p);
      console.log(response);
    });
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
  function getInfo(bearer, deviceName, next) {
    let time1 = performance.now();
    let settings = {
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
  function createDevice(name, func, nextCall) {
    let baseName = "timeDisconnect1";
    let time1 = performance.now();
    // let time2 = performance.now();
    let numDevices = 1;
    for (let i = 0; i < numDevices; i++) {
      let deviceName = "updateThreeTime2";
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
