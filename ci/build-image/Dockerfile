# Based on: https://github.com/docker-library/golang/tree/master/1.14/alpine3.11
# But switched out alpine for docker as base image
FROM docker

RUN apk add --no-cache ca-certificates go

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
