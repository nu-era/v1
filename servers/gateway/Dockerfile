FROM scratch
COPY gateway /gateway
ENTRYPOINT [ "/gateway" ]
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true