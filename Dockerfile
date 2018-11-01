FROM golang:alpine

ADD . /go/src/github.com/wheelerlaw/octodns
WORKDIR /go/src/github.com/wheelerlaw/octodns
RUN apk update && apk add make git zip && make

FROM alpine:latest
COPY --from=0 /go/src/github.com/wheelerlaw/octodns/octodns /usr/local/bin

CMD "octodns"
