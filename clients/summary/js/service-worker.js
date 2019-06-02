var device = "";

// gets distance between two points in KM
function distance(lat1, lng1, lat2, lng2, miles) { // miles optional
    if (typeof miles === "undefined"){miles=false;}
    function deg2rad(deg){return deg * (Math.PI/180);}
    function square(x){return Math.pow(x, 2);}
    var r=6371; // radius of the earth in km
    lat1=deg2rad(lat1);
    lat2=deg2rad(lat2);
    var lat_dif=lat2-lat1;
    var lng_dif=deg2rad(lng2-lng1);
    var a=square(Math.sin(lat_dif/2))+Math.cos(lat1)*Math.cos(lat2)*square(Math.sin(lng_dif/2));
    var d=2*r*Math.asin(Math.sqrt(a));
    if (miles){
        return d * 0.621371;
        } //return miles
    else{
        return d;
    } //return km
}

var circleData = {
    4:  {
        'message': 'Light Shaking Expected',
        'color': '#7efbdf'
    },
    5: { 
        'message': 'Moderate Shaking Expected',
        'color': '#95f879'
    },
    6: {
        'message': 'Strong Shaking Expected',
        'color': '#f7f835'
    },
    7: {
        'message': 'Very Strong Shaking Expected',
        'color': '#fdca2c'
    },
    8: {
        'message': 'Severe Shaking Expected',
        'color': '#ff701f'
    },
    9: {
        'message': 'Violent Shaking Expected',
        'color': '#ec2516'
    },
    10: {
        'message': 'Extreme Shaking Expected',
        'color': '#c81e11'
    }   
}

function getTime(lat, long, m, res) {
    let dist;
    if (lat !== null && long !== null) {
        dist = distance(parseFloat(res[0]), parseFloat(res[1]), lat, long)
    }
    // number of seconds using speed of 3km/s
    let tmp = dist / 3
    date = new Date(m.orig_time)
    date.setSeconds(date.getSeconds() + tmp)
    return date;
}


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