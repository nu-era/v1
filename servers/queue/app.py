from flask import Flask, request, Response
import json
from datetime import datetime, timezone
import config, pika, sys, stomp, xmltodict, ssl, pymongo, ast
from shapely.geometry import Point
from shapely.geometry.polygon import Polygon
import geopy.distance

flask_app = Flask(__name__)
flask_app.app_context().push()




# RabbitMQ declaration
creds = pika.PlainCredentials(config.rUSER, config.rPW)
conn = pika.BlockingConnection(pika.ConnectionParameters(host=config.mqHOST, port=config.mqPORT, credentials=creds, heartbeat=0))
mq_chan = conn.channel()
mq_chan.queue_declare(queue=config.qName, durable=True)




# what to do on message from ShakeAlert
class MyListener(stomp.ConnectionListener):
    def on_message(self, headers, message):
        print('received a message "%s", headers: %s' % (type(message),headers))
        # have contour message to send out
        if int(headers['subscription']) == 3:
            event = process(message)
            event = makePolygons(event)
            print(event)
            mq_chan.basic_publish(exchange='', routing_key=config.qName, body=json.dumps(event))
        



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
            event['MMI_' + c['MMI']['#text']] = c['polygon']['#text']

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


# def filterUsers()


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
    "Error subscribing to topic"
    # stop connections
    sa_conn.disconnect()
    test_conn.disconnect()




# run app
if __name__ == "__main__":
    flask_app.run(debug=False, host='msg', port=5000)
    # host=config.host