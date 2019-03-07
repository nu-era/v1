# publish.py
import pika, os, datetime

# Access the CLODUAMQP_URL environment variable and parse it (fallback to localhost)
url = os.environ.get('CLOUDAMQP_URL', 'amqp://rkbyzjjk:1hAf-c13GnTylFcLNyphTtI6k8qN3031@beaver.rmq.cloudamqp.com/rkbyzjjk')
params = pika.URLParameters(url)
connection = pika.BlockingConnection(params)
channel = connection.channel() # start a channel
channel.queue_declare(queue='hello') # Declare a queue
channel.basic_publish(exchange='',
                      routing_key='hello',
                      body='Nu-era test message for Marlina!')

print(" [x] Sent 'Hello World!'")
connection.close()
