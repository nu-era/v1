const device;

self.addEventListener('message', function(event){
    device = JSON.parse(event.data);
});

self.addEventListener("push", e => {
    const data = e.data.json();
    let loc = data.location.split(',')
    // get time equake will hit user
    let d = getTime(device.latitude, device.longitude, data, loc)
    let curr = new Date();
    
    self.registration.showNotification("EARTHQUAKE ALERT!", {
    body: "Expected Intensity is: " + 
            data.intensity + " <br/> " + 
            circleData[data.intensity].message + " <br/> Estimated Time to Impact: " + 
            ((d.getTime() - curr.getTime()) / 1000) + " seconds",
    vibrate: [300, 100, 400, 300, 100, 400, 300, 100, 400, 300, 100, 400, 300, 100, 400]
    });
});