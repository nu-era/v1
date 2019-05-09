(function() {
  "use strict";
  window.addEventListener("load", init);

  function init() {
    let baseName = "stressTestTwo";
    let time1 = performance.now();
    let time2 = performance.now();
    for (let i = 0; i < 1; i++) {
      let deviceName = baseName + i;
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
      $.ajax(settings).done(function (response) {
        time2 = performance.now();
        console.log(response);
      });
    }
    // let time2 = performance.now();
    let p = document.getElementById("time-info");
    p.innerHTML = "Creating 100 devices took " + ((time2 - time1)/1000) + " seconds.";

  }
})();
