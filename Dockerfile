FROM golang:1.18-alpine AS build-env
ENV GOPROXY=goproxy.cn,direct
ADD . /go/src/app
WORKDIR /go/src/app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --update add git curl tzdata && \
    go build -v -o /go/src/app/kas main.go && \
    curl https://archive.org/download/kindlegen2.9/kindlegen_linux_2.6_i386_v2_9.tar.gz | tar -zx

FROM alpine
COPY --from=build-env /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build-env /go/src/app/kas /app/kas
COPY --from=build-env /go/src/app/kindlegen /bin/kindlegen
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --update add --no-cache curl
WORKDIR /app
VOLUME ["/app/storage"]
HEALTHCHECK --interval=1m --timeout=10s \
  CMD curl -f http://localhost:1323/ping || exit 1
EXPOSE 1323
cmd ./kas
