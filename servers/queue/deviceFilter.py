from shapely.geometry import Polygon, Point
from pymongo import MongoClient
from pprint import pprint


args = sys.argv #gets arguments from script call
if len(args) < 2:
    print("error: please pass in polygon information")
else:
    poly = args[1] # Gets the polygon object

    # Connects to mongoDB
    client = MongoClient(port=27017)
    result = []
    db = client.db
    collection = db.devices

    # Gets all devices
    devices = collection.find()
    for device in devices:
        location = Point(device["Lat"], device["Long"]) # Create point for device
        onPolygon = poly.touches(location) # Check if device is on edge of polygon
        inPolygon = poly.contains(location) # check if device is inside of polygon
        if inPolygon or onPolygon:
            result.append(device)

    return result
