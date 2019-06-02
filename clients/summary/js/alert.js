// build Google API map with epicenter of earthquake as center 
function initMap(lat, long, msg) {
    if ((lat == null) || (long == null)) {
        lat = parseFloat(-34.397)
        long = parseFloat(150.644)
    }
    bounds = new google.maps.LatLngBounds();
    map = new google.maps.Map(document.getElementById('map'), {
        center: {lat: lat, lng: long},
        zoom: 11

    });

    marker = new google.maps.Marker({
        position: map.center,
        map: map,
        title: 'Epicenter of earthquake',
        label: "Epicenter"
    });
    bounds.extend(marker.position)
}


// creates custom timestamp
function timeStamp() {
    // Create a date object with the current time
      var now = new Date();

    // Create an array with the current month, day and time
      var date = [ now.getMonth() + 1, now.getDate(), now.getFullYear() ];

    // Create an array with the current hour, minute and second
      var time = [ now.getHours(), now.getMinutes(), now.getSeconds(), now.getMilliseconds()];

    // Determine AM or PM suffix based on the hour
      var suffix = ( time[0] < 12 ) ? "AM" : "PM";

    // Convert hour from military time
      time[0] = ( time[0] < 12 ) ? time[0] : time[0] - 12;

    // If hour is 0, set it to 12
      time[0] = time[0] || 12;

    // If seconds and minutes are less than 10, add a zero
      for ( var i = 1; i < 3; i++ ) {
        if ( time[i] < 10 ) {
          time[i] = "0" + time[i];
        }
      }

    // Return the formatted string
      return date.join("/") + " " + time.join(":") + " " + suffix;
  }


  // check if a passed string contains JSON
  function IsJsonString(str) {
    try {
        JSON.parse(str);
    } catch (e) {
        return false;
    }
    return true;
}

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


function addCircle(m, circleData, map) {
    var circle = new google.maps.Circle({
        map: map,
        radius: m.radius,
        fillColor: circleData[m.intensity].color,
        center: map.center
    })
    return circle;
}