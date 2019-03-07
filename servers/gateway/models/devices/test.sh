docker run -d \
-p 27017:27017 \
--name mongo_test \
mongo

go test

docker rm -f mongo_test