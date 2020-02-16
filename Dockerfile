FROM golang:alpine AS build-env
ENV GO111MODULE=on
ADD . /go/src/app
WORKDIR /go/src/app
RUN apk --update add git tzdata && \
    go build -v -o /go/src/app/hcc main.go && \
    export GO111MODULE=off && \
    go get github.com/GeertJohan/go.rice && \
    go get github.com/GeertJohan/go.rice/rice && \
    rice append --exec /go/src/app/hcc

FROM alpine
COPY --from=build-env /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build-env /go/src/app/hcc /app/hcc
WORKDIR /app
VOLUME ["/app/storage"]
EXPOSE 1323
cmd ./hcc