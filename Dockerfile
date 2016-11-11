FROM golang:latest
MAINTAINER <284077318@qq.com>

COPY . go/src/github.com/jicg/easypos

WORKDIR go/src/github.com/jicg/easypos

VOLUME go/src/github.com/jicg/easypos/data
VOLUME go/src/github.com/jicg/easypos/log
RUN go get github.com/jicg/easypos
RUN go install -a github.com/jicg/easypos

EXPOSE 4000
CMD easypos web --port 4000