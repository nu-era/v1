#!/usr/bin/env bash
set -e
./build.sh
docker push newera/gateway
docker push newera/mysql  

#export TLSCERT=/etc/letsencrypt/live/api.ericjcwei.me/fullchain.pem
#export TLSKEY=/etc/letsencrypt/live/api.ericjcwei.me/privkey.pem
export MESSAGESADDR="messaging:80"
export SUMMARYADDR="summary:80"
export REDISADDR="redisserver:6379"

ssh -i ~/.ssh/MyPrivKey.pem ec2-54-68-59-121.us-west-2.compute.amazonaws.com 'bash -s' << EOF
#Cleanup existing docker images
docker pull newera/gateway
docker pull newera/mysql
docker rm -f gateway
docker rm -f mysql
docker rm -f redisserver

# Create docker network
docker network create apinet

# Run redis instance
docker run -d \
--network apinet \
--name redisserver \
redis

# Run mysql instance
docker run -d \
--network apinet \
--name mysql \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
-e MYSQL_DATABASE=mysql \
ericjwei/mysql

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
-e MESSAGESADDR=$MESSAGESADDR \
-e SUMMARYADDR=$SUMMARYADDR \
-e MYSQL_ROOT_PASSWORD=$MYSQL_ROOT_PASSWORD \
ericjwei/gateway

docker system prune -f
docker volume prune -f

exit

EOF