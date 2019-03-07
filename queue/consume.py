# consume.py
import pika, os

# Access the CLODUAMQP_URL environment variable and parse it (fallback to localhost)
url = os.environ.get('CLOUDAMQP_URL', 'amqp://rkbyzjjk:1hAf-c13GnTylFcLNyphTtI6k8qN3031@beaver.rmq.cloudamqp.com/rkbyzjjk')
params = pika.URLParameters(url)
connection = pika.BlockingConnection(params)
channel = connection.channel() # start a channel

channel.queue_declare(queue='hello')

def callback(ch, method, properties, body):
  print(" [x] Received %r" % body)

channel.basic_consume(callback,
                      queue='hello',
                      no_ack=True)

print(' [*] Waiting for messages:')
channel.start_consuming()
