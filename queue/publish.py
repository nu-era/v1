# publish.py
#import pika, os, datetime

# Access Docker container running RabbitMQ
# creds = pika.PlainCredentials(config.rUSER, config.rPW)
# conn = pika.BlockingConnection(pika.ConnectionParameters(host=config.mqHOST, port=config.mqPORT, credentials=creds, heartbeat=0))
# mq_chan = conn.channel()
# channel.queue_declare(queue='alerts') # Declare a queue
# channel.basic_publish(exchange='',
#                       routing_key='alerts',
#                       body='Nu-era test message for Marlina!')

# print(" [x] Sent 'Hello World!'")
# connection.close()

import pika
import sys, os
import json
import datetime, time

url = os.environ.get('CLOUDAMQP_URL', 'http://ec2-54-68-59-121.us-west-2.compute.amazonaws.com')
params = pika.URLParameters(url)
connection = pika.BlockingConnection(params)
channel = connection.channel()

channel.exchange_declare(exchange='logs',
                         exchange_type='fanout')

data = {}
data['message'] = 'hello'
data['status'] = 'success'
ts = time.time()
data['time_stamp'] = datetime.datetime.fromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')
data['location'] = 'vancouver'
json_data = json.dumps(data)
message = json_data
channel.basic_publish(exchange='logs',
                      routing_key='',
                      body=message)
print(" [x] Sent %r" % message)
connection.close()
