#!/usr/bin/env bash
./build.sh
export TLSCERT="/tls/fullchain.pem"
export TLSKEY="/tls/privkey.pem"
# export MESSAGESADDR="messaging:80"
# export SUMMARYADDR="summary:80"

docker rm -f mongodb
docker run -d \
-p 27017:27017 \
--network apinet \
--name mongodb \

docker rm -f test_gateway
docker run -d \
--name test_gateway \
--network apinet \
-p 443:443 \
-v "$(pwd)"/tls:/tls:ro \
-e TLSCERT=$TLSCERT \
-e TLSKEY=$TLSKEY \
# -e MESSAGESADDR=$MESSAGESADDR \
# -e SUMMARYADDR=$SUMMARYADDR \
# -e SESSIONKEY=$SESSIONKEY \
ericjwei/gateway
docker system prune -f