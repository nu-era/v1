FROM scratch
COPY gateway /gateway
ENTRYPOINT [ "/gateway" ]
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*