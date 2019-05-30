#!/usr/bin/env bash
set -e
./build.sh
./db/deploy.sh
docker push bfranzen1/newera-gateway
#docker push newera/mysql  

# build goqueue microservice container
(cd ../goqueue/ ; sh build.sh)
docker push bfranzen1/goqueue

source ./twilio.env
export TLSCERT=/etc/letsencrypt/live/api.bfranzen.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.bfranzen.me/privkey.pem
export MYSQL_ROOT_PASSWORD="shakealert"
export REDISADDR="redisserver:6379"
export MONGOADDR="mgo:27017"
export SESSIONKEY="shakealert"
export RABBITMQ="rmq:5672"
export GOQ="queue:5000"
#export WCADDRS="wc:8000"

echo "Connecting to server..."
ssh ec2-user@ec2-34-212-199-173.us-west-2.compute.amazonaws.com 'bash -s' << EOF
#Clean up existing docker images
printf 'y' | docker system prune -a --volumes;

# Create docker network
docker network create apinet;

docker pull bfranzen1/goqueue
docker pull bfranzen1/newera-gateway
# docker pull newera-mysql
docker rm -f gateway
# docker rm -f mysql
docker rm -f redisserver
docker rm -f mgo
docker rm -f rmq

# Run RabbitMQ instance
docker run -d \
--name rmq \
--network apinet \
-p 5672:5672 \
-p 15672:15672 \
rabbitmq;

# Run mongo instance
docker run -d \
-p 27017:27017 \
--network apinet \
--name mgo \
mongo;

# Run redis instance
docker run -d \
--network apinet \
--name redisserver \
redis;


# Run mysql instance
# docker rm -f AlertsDB
# docker run -d --name AlertsDB --network apinet \
# -e MYSQL_ROOT_PASSWORD=\$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=AlertsDB \
# alerts;

# docker run -d \
# --network apinet \
# --name mysql \
# -e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
# -e MYSQL_DATABASE=mysql \
# newera/mysql;

sleep 10s; # need to wait for rmq for some reason

# Run web server
docker run -d \
-p 443:443 \
--name gateway \
-v /etc/letsencrypt:/etc/letsencrypt:ro \
--network apinet \
-e SESSIONKEY=$SESSIONKEY \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
-e REDISADDR=$REDISADDR \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MONGO_ADDR=$MONGOADDR \
-e RABBITMQ=$RABBITMQ \
-e GOQ=$GOQ \
-e TWILIO_ACCOUNT_SID=$TWILIO_ACCOUNT_SID \
-e TWILIO_AUTH_TOKEN=$TWILIO_AUTH_TOKEN \
bfranzen1/newera-gateway;
#-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \

exit

EOF