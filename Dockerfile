FROM golang:alpine AS build-env
ENV GO111MODULE=on
ADD . /go/src/app
RUN  apk --update add git tzdata
WORKDIR /go/src/app
RUN go build -v -o /go/src/app/hcc cmd/main.go

FROM alpine
COPY --from=build-env /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
COPY --from=build-env /go/src/app/hcc /app/hcc
WORKDIR /app
VOLUME ["/app/storage"]
EXPOSE 1323
cmd ./hcc