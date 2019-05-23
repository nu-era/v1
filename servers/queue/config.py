import os
# environment variables
# host = os.environ["ADDR"]
# mqHOST = os.environ["RABBITMQ_HOST"]
# mqPORT = os.environ["RABBITMQ_PORT"]
# rUSER = os.environ["RABBITMQ_USER"] 
# rPW = os.environ["RABBITMQ_PW"]
# qName = os.environ["RMQUEUE"]
DM_USER = os.environ["DM_USER"]
DM_PW = os.environ["DM_PW"]
STOMP_PORT = os.environ["STOMP_PORT"]
amq_broker = os.environ["AMQ_BROKER"]
test_host = os.environ["TEST_BROKER"]


# topics to subscribe to for ShakeAlert 
gmcontour_topic = "/topic/eew.sys.gm-contour.data"
heartbeat_topic = "/topic/eew.sys.ha.data"

# test service topic
contour_test = "/topic/eew.test_ericjwei.gm-contour.data"