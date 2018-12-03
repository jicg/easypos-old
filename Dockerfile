FROM golang:latest as builder
MAINTAINER <284077318@qq.com>
COPY . $GOPATH/src/github.com/jicg/easypos
WORKDIR $GOPATH/src/github.com/jicg/easypos
WORKDIR $GOPATH/src/github.com/jicg/easypos
RUN go get  github.com/jicg/easypos
RUN go install -a github.com/jicg/easypos

FROM scratch AS final
MAINTAINER <284077318@qq.com>
COPY --from=builder /go/bin/easypos /app/easypos
COPY --from=builder /go/src/github.com/jicg/easypos/views /app/views
COPY --from=builder /go/src/github.com/jicg/easypos/public /app/public
VOLUME /app/data
VOLUME /app/log
EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["/app/easypos"]
