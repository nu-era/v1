from flask import Flask, request, Response
import json
from datetime import datetime, timezone
import config, pika, sys, stomp, xmltodict, ssl, ast
from pymongo import MongoClient
from shapely.geometry import Point
from shapely.geometry.polygon import Polygon
from bson.json_util import dumps
import geopy.distance

flask_app = Flask(__name__)
print("WOW I WORK HERE CUZ IM DUMB")
flask_app.app_context().push()


# RabbitMQ declaration
creds = pika.PlainCredentials(config.rUSER, config.rPW)
conn = pika.BlockingConnection(pika.ConnectionParameters(host=config.mqHOST, port=config.mqPORT, credentials=creds, heartbeat=0))
mq_chan = conn.channel()
mq_chan.queue_declare(queue=config.qName, durable=True)


# Connects to mongoDB
client = MongoClient(config.mgo_host, int(config.mgo_port))
db = client.db
collection = db.devices




# what to do on message from ShakeAlert
class MyListener(stomp.ConnectionListener):
    def on_message(self, headers, message):
        # have contour message to send out
        if int(headers['subscription']) == 3:
            event = process(message)
            event = makePolygons(event)
            pushToUsers(event)
        



# takes in a message and converts to python dictionary to parse
# and extract data
def process(message):
    # convert to dict
    d = xmltodict.parse(message)
    contours = d['event_message']['gm_info']['gmcontour_pred']['contour']
    
    # build event to push onto message queue
    event = {}
    event['magnitude'] = d['event_message']['core_info']['mag']['#text']
    event['location'] = d['event_message']['core_info']['lat']['#text'] + ', ' + d['event_message']['core_info']['lon']['#text']
    event['orig_time'] = d['event_message']['core_info']['orig_time']['#text']
    event['orig_time_unc'] = d['event_message']['core_info']['orig_time_uncer']['#text']

    # add affected areas by severity 
    for c in contours:
        if float(c['MMI']['#text']) >= 4:
            event['MMI_' + c['MMI']['#text'][0]] = c['polygon']['#text']

    return event

        
# function that takes in a parsed event and returns areas affected by different intensities
def makePolygons(event):
    locations = {k: v for k, v in event.items() if k.startswith('MMI_')}
    polygons = {}
    for key, value in locations.items():
        # change event string to polygon points
        points = value.split(' ')
        points = ''.join(['(' + l + ')' for l in points])
        points = points.replace(')(', '), (')
        
        # add polygon to list
        polygon = Polygon(ast.literal_eval(points))
        polygons[key] = polygon

        # add radius in meters for area affected
        # gets distance between top point and epicenter
        event[key + "_radius"] = geopy.distance.distance(polygon.exterior.coords[0], event['location']).m

    # add polygons to event dict
    event["areas_affected"] = polygons

    # return updated event dict
    return event




# filters devices by checking if they are within/touching the passed polygon/area affected
def filterDevices(polygon):
    devices = []
    deviceIDs = []
    # Gets all devices
    devices = collection.find()

    for device in devices:
        location = Point(device["Lat"], device["Long"]) # Create point for device

        onPolygon = polygon.touches(location) # Check if device is on edge of polygon
        inPolygon = polygon.contains(location) # check if device is inside of polygon
        print("LAT: ", device["Lat"], file=sys.stderr)
        app.logger.info("LONG: ", device["Long"])
        print("LOCATION: ", noLoc)
        noLoc = (device["Lat"] == None and device["Long"] == None) # device didn't provide location, notify anyway
        sys.stdout.flush()
        if inPolygon or onPolygon or noLoc:
            devices.append(device)
            deviceIDs.append(device['ID'])
        

    return devices, deviceIDs




# function that pushes alerts to message queue to notify subset of users
def pushToUsers(event):
    # get list of users starting with highest MMI intensity (10 -> 4)
    for x in range(10, 3, -1):
        if 'MMI_' + str(x) in event:
            u_event = {
                'magnitude': event['magnitude'],
                'location': event['location'],
                'orig_time': event['orig_time'],
                'orig_time': event['orig_time'],
                'intensity': x,
                'radius': event['MMI_' + str(x) + '_radius']
            }
            devices, u_event['deviceIDs'] = filterDevices(event['areas_affected']['MMI_' + str(x)])
            u_event['devices'] = dumps(devices)
            mq_chan.basic_publish(exchange='', routing_key=config.qName, body=json.dumps(u_event))




# connects via stomp to server hosted at host on port using the passed credentials
def makeConnection(user, pw, host, port):
    conn = stomp.Connection([(host,port)], auto_decode=True)
    conn.set_listener('', MyListener())
    conn.set_ssl(for_hosts=[(host,port)],ssl_version=ssl.PROTOCOL_TLS)
    conn.start()
    conn.connect(user, pw, wait=True)
    return conn




# connect to ShakeAlert
sa_conn = makeConnection(config.DM_USER, config.DM_PW, config.amq_broker, config.STOMP_PORT)

# connect to Test Service to demo/run through test scenarios
test_conn = makeConnection(config.DM_USER, config.DM_PW, config.test_host, config.STOMP_PORT)




# handle subscriptions
try:
    sa_conn.subscribe(destination=config.heartbeat_topic, id=1, ack='auto')
    sa_conn.subscribe(destination=config.gmcontour_topic, id=3, ack='auto')
    test_conn.subscribe(destination=config.contour_test, id=3, ack='auto')
except Exception as e:
    print("Error subscribing to topic", file=sys.stderr)
    # stop connections
    sa_conn.disconnect()
    test_conn.disconnect()




# run app
if __name__ == "__main__":
    print("STARTING APP")
    print(dumps(collection.find()))
    flask_app.run(debug=False, host=config.host, port=5000)