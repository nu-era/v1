FROM python:3.6-alpine
WORKDIR /app
COPY . /app
CMD ["python","-u","app.py"]
# dependencies for shapely 
RUN apk --update add build-base libxslt-dev

RUN apk add --virtual .build-deps \
        --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing \
        --repository http://dl-cdn.alpinelinux.org/alpine/edge/main \
        gcc libc-dev geos-dev geos && \
    runDeps="$(scanelf --needed --nobanner --recursive /usr/local \
    | awk '{ gsub(/,/, "\nso:", $2); print "so:" $2 }' \
    | xargs -r apk info --installed \
    | sort -u)" && \
    apk add --virtual .rundeps $runDeps

RUN geos-config --cflags

RUN pip install --disable-pip-version-check -r requirements.txt

RUN apk del build-base python3-dev && \
    rm -rf /var/cache/apk/*

# RUN echo "http://mirror.leaseweb.com/alpine/edge/testing/" >> /etc/apk/repositories
# RUN apk add --no-cache geos-dev
#RUN pip3 install -r requirements.txt
EXPOSE 5000-5110
ENTRYPOINT ["python3", "app.py"]