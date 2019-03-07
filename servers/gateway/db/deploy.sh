# upload newest container
docker build -t bfranzen1/alerts .
docker login
docker push bfranzen1/alerts 

# run container

ssh ec2-user@ec2-54-218-119-10.us-west-2.compute.amazonaws.com "
docker volume rm $(docker volume ls -qf dangling=true);

docker rm -f AlertsDB
docker run -d --name AlertsDB --network api \
-e MYSQL_ROOT_PASSWORD=\$MYSQL_ROOT_PASSWORD -e MYSQL_DATABASE=AlertsDB \
bfranzen1/alerts;
"
