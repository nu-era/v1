# Queue Microservice
FROM alpine 
RUN apk add --no-cache ca-certificates
COPY goqueue /goqueue
EXPOSE 8000
ENTRYPOINT [ "/goqueue" ]