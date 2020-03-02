FROM golang:alpine AS build-env
ADD . /go/src/app
WORKDIR /go/src/app
RUN apk --update add git curl tzdata && \
    export GO111MODULE=off && \
    go get github.com/GeertJohan/go.rice && \
    go get github.com/GeertJohan/go.rice/rice && \
    rice embed-go && \
    export GO111MODULE=on && \
    go build -v -o /go/src/app/kas main.go && \
    curl http://kindlegen.s3.amazonaws.com/kindlegen_linux_2.6_i386_v2_9.tar.gz | tar -zx

FROM alpine
COPY --from=build-env /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build-env /go/src/app/kas /app/kas
COPY --from=build-env /go/src/app/kindlegen /bin/kindlegen
WORKDIR /app
VOLUME ["/app/storage"]
EXPOSE 1323
cmd ./kas