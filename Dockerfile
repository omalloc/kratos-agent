FROM docker.bs58i.baishancloud.com/base/alpine:3.14

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
     && apk add --no-cache tzdata

COPY /bin /app

WORKDIR /app

EXPOSE 8000
EXPOSE 9000
VOLUME /data/conf

ENV TZ=Asia/Shanghai

CMD ["./server", "-conf", "/data/conf"]
