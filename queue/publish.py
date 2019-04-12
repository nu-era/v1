
#!/usr/bin/env python
import pika
import sys, os
import json
import datetime, time
import geocoder

shaking = {
    4:"Light",
    5:"Moderate",
    6:"Strong",
    7:"Very Strong",
    8:"Severe",
    9:"Violent",
    10:"Extreme"
}

message_structure = "Earthquake warning. {0} shaking expected in an estimated {1} seconds. "
message = ""
information = {}
def send():
    global message, information
    url = os.environ.get('CLOUDAMQP_URL', 'http://ec2-54-68-59-121.us-west-2.compute.amazonaws.com')
    params = pika.URLParameters(url)
    connection = pika.BlockingConnection(params)
    channel = connection.channel()

    channel.exchange_declare(exchange='logs',
                             exchange_type='fanout')
    g = geocoder.ip('me')
    data = {}
    print(message)
    data['message'] = message
    ts = time.time()
    data['time_stamp'] = datetime.datetime.fromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')
    data['location'] = information["Location"]
    data['magnitude'] = information["Magnitude"]
    data['intensity'] = information["Intensity"]
    data['countdown'] = information["Time"]
    data['Devices'] = information["DeviceIDs"]
    json_data = json.dumps(data)
    message = json_data
    channel.basic_publish(exchange='logs',
                          routing_key='',
                          body=message,
                          properties=pika.BasicProperties(
                         delivery_mode = 2, # make message persistent
                      ))
    print(" [x] Sent %r" % message)
    connection.close()

def process_data():
    global information
    information["Location"] = "location"
    information["Magnitude"] = 3.2
    information["Intensity"] = 5
    information["Time"] = 35
    information["DeviceIDs"] = ["43287", '43829', '57290', '50134']

def construct_message():
    global information, shaking, message, message_structure
    intensity = information.get("Intensity")
    time = information.get("Time")
    shaking = shaking.get(intensity)
    message = message_structure.format(shaking, time)

# def get_information(location, magnitude, intensity, affected_area):
def main():
    process_data()
    construct_message()
    send()

if __name__ == '__main__':
    main()
