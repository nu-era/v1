#!/usr/bin/env bash
set -e
./build.sh
docker push bfranzen1/newera-gateway
#docker push newera/mysql  

export TLSCERT=/etc/letsencrypt/live/api.bfranzen.me/fullchain.pem
export TLSKEY=/etc/letsencrypt/live/api.bfranzen.me/privkey.pem
export REDISADDR="redisserver:6379"
export MONGOADDR="mgo:27017"
export SESSIONKEY=$SESSIONKEY

echo "Connecting to server..."
ssh ec2-user@ec2-34-212-199-173.us-west-2.compute.amazonaws.com 'bash -s' << EOF
#Cleanup existing docker images
docker pull bfranzen1/newera-gateway
#docker pull newera-mysql
docker rm -f gateway
#docker rm -f mysql
docker rm -f redisserver
docker rm -f mgo

# Create docker network
#docker network create apinet

# Run mongo instance
docker run -d \
-p 27017:27017 \
--network apinet \
--name mgo \
mongo

# Run redis instance
docker run -d \
--network apinet \
--name redisserver \
redis

# Run mysql instance
#docker run -d \
#--network apinet \
#--name mysql \
#-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
#-e MYSQL_DATABASE=mysql \
#newera/mysql

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
-e MONGO_ADDR=$MONGOADDR \
bfranzen1/newera-gateway

#-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \

docker system prune -f
docker volume prune -f

exit

EOF