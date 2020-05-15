# Based on: https://github.com/docker-library/golang/tree/master/1.14/alpine3.11
# But switched out alpine for docker as base image
FROM docker

RUN apk add --no-cache ca-certificates go git make protobuf

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

# JUnit report
RUN go get -v -u github.com/jstemmer/go-junit-report

# Controller-gen
RUN go get -v -u sigs.k8s.io/controller-tools/cmd/controller-gen

# Mockgen
RUN GO111MODULE=on go get -v -u github.com/golang/mock/mockgen@latest

# Protobuf
RUN GO111MODULE=on go get github.com/golang/protobuf@v1.4.0
