# Build Docker Container
docker build -t bfranzen1/queue .

docker login
docker push bfranzen1/queue