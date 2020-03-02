FROM golang:alpine AS build-env
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
ADD . /go/src/app
WORKDIR /go/src/app
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --update add git curl tzdata && \
    go build -v -o /go/src/app/kas main.go && \
    echo "done!" && \
    export GO111MODULE=off && \
    go get github.com/GeertJohan/go.rice && \
    go get github.com/GeertJohan/go.rice/rice && \
    rice append --exec kas && \
    curl http://kindlegen.s3.amazonaws.com/kindlegen_linux_2.6_i386_v2_9.tar.gz | tar -zx

FROM alpine
COPY --from=build-env /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build-env /go/src/app/kas /app/kas
COPY --from=build-env /go/src/app/public /app/public
COPY --from=build-env /go/src/app/kindlegen /bin/kindlegen
WORKDIR /app
VOLUME ["/app/storage"]
EXPOSE 1323
cmd ./kas