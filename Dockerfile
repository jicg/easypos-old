FROM golang:latest as builder
MAINTAINER <284077318@qq.com>
COPY . $GOPATH/src/github.com/jicg/easypos
WORKDIR $GOPATH/src/github.com/jicg/easypos
WORKDIR $GOPATH/src/github.com/jicg/easypos
RUN go get  github.com/jicg/easypos
RUN go install -a github.com/jicg/easypos

FROM debian:latest
#daocloud.io/centos:latest
#alpine:latest
MAINTAINER <284077318@qq.com>
ADD go/bin/easypos /usr/bin/easypos

ADD go/src/github.com/jicg/easypos/views /app/views
ADD go/src/github.com/jicg/easypos/public /app/public
VOLUME /app/data
VOLUME /app/log
EXPOSE 4000
WORKDIR /app
CMD /usr/bin/easypos web --port 4000
