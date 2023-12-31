FROM alpine:3.15

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
     && apk add --no-cache tzdata curl

COPY ./bin/ /app
COPY ./configs /data/conf

WORKDIR /app

EXPOSE 2381
EXPOSE 2382
VOLUME /data/conf

ENV TZ=Asia/Shanghai

HEALTHCHECK --interval=5s --timeout=5s --start-period=3s --retries=3 \
    CMD curl -sS 'http://127.0.0.1:2381/health' || exit 1

CMD ["./server", "--conf", "/data/conf"]
