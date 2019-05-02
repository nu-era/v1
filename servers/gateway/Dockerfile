FROM alpine
RUN apk add --no-cache ca-certificates
COPY gateway /gateway
ENTRYPOINT [ "/gateway" ]
