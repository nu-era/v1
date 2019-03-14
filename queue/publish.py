
#!/usr/bin/env python
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
