from flask import Flask, request, Response
import json
from models import Member, Message, Channel, db
from datetime import datetime, timezone
import config
import pika
import sys

flask_app = Flask(__name__)
flask_app.app_context().push()

# RabbitMQ declaration
creds = pika.PlainCredentials(config.rUSER, config.rPW)
conn = pika.BlockingConnection(pika.ConnectionParameters(host=config.mqHOST, port=config.mqPORT, credentials=creds, heartbeat=0))
mq_chan = conn.channel()
mq_chan.queue_declare(queue=config.rmQueue, durable=True)

if __name__ == "__main__":
    flask_app.run(debug=False, host="msg", port=5000)
