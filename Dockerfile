FROM golang:latest

MAINTAINER <284077318@qq.com>

COPY . $GOPATH/src/github.com/jicg/easypos

WORKDIR $GOPATH/src/github.com/jicg/easypos

WORKDIR $GOPATH/src/github.com/jicg/easypos

RUN go get github.com/jicg/easypos

RUN go install -a github.com/jicg/easypos

# EXPOSE 4000
# CMD easypos web --port 4000