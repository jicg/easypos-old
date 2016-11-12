FROM golang:latest

MAINTAINER <284077318@qq.com>

COPY . $GOPATH/src/github.com/jicg/easypos

WORKDIR $GOPATH/src/github.com/jicg/easypos

WORKDIR $GOPATH/src/github.com/jicg/easypos

RUN go get -tags netgo -installsuffix netgo github.com/jicg/easypos  

RUN go install -a -tags netgo -installsuffix netgo github.com/jicg/easypos

# EXPOSE 4000
# CMD easypos web --port 4000