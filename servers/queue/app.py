from flask import Flask, request, Response
import json
from models import Member, Message, Channel, db
from datetime import datetime, timezone
import config
import pika
import sys
import stomp
import xmltodict


flask_app = Flask(__name__)
flask_app.app_context().push()

# RabbitMQ declaration
creds = pika.PlainCredentials(config.rUSER, config.rPW)
conn = pika.BlockingConnection(pika.ConnectionParameters(host=config.mqHOST, port=config.mqPORT, credentials=creds, heartbeat=0))
mq_chan = conn.channel()
mq_chan.queue_declare(queue=config.rmQueue, durable=True)


# what to do on message from ShakeAlert
class MyListener(stomp.ConnectionListener):
    def on_connected(self, headers, body):
        print('Connected: "%s"' % headers)
    def on_error(self, headers, message):
        print('received an error "%s"' % message)
    def on_message(self, headers, message):
        print('received a message "%s", headers: %s' % (type(message),headers))
        process((message))

DM_USER = os.environ["DM_USER"]
DM_PW = os.environ["DM_PW"]
STOMP_PORT = os.environ["STOMP_PORT"]
amq_broker = os.environ["AMQ_BROKER"]

# connect to ShakeAlert
sa_conn = stomp.Connection([(amq_broker,STOMP_PORT)],auto_decode=True)
sa_conn.set_listener('', MyListener())
sa_conn.set_ssl(for_hosts=[(amq_broker,STOMP_PORT)],ssl_version=ssl.PROTOCOL_TLS)
sa_conn.start()
sa_conn.connect(DM_USER, DM_PW, wait=True)

# topics to subscribe to
gmcontour_topic = "/topic/eew.sys.gm-contour.data"
heartbeat_topic = "/topic/eew.sys.ha.data"

# handle subscriptions
try:
    sa_conn.subscribe(destination=heartbeat_topic, id=1, ack='auto')
    sa_conn.subscribe(destination=gmcontour_topic, id=3, ack='auto')
except Exception as e:
    "Error subscribing to topic {}: {}".format(topic,e)
    sa_conn.disconnect()
    sys.exit()


if __name__ == "__main__":
    flask_app.run(debug=False, host="msg", port=5000)
