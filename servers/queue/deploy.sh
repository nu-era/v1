

ssh ec2-user@ec2-54-68-59-121.us-west-2.compute.amazonaws.com "
docker rm -f rmq
docker run -d --name rmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
"