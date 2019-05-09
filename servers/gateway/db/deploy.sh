# upload newest container
docker build -t bfranzen1/alerts .
docker push bfranzen1/alerts 
export MYSQL_ROOT_PASSWORD="shakealert"
# run container

ssh ec2-user@ec2-34-212-199-173.us-west-2.compute.amazonaws.com "
docker volume rm $(docker volume ls -qf dangling=true);
docker login

docker rm -f alerts
docker run -d --name alerts --network apinet \
-e MYSQL_ROOT_PASSWORD=shakealert -e MYSQL_DATABASE=AlertsDB \
bfranzen1/alerts;
"
