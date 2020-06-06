# Based on: https://github.com/docker-library/golang/tree/master/1.14/alpine3.11
# But switched out alpine for docker as base image
FROM docker

RUN apk add --no-cache ca-certificates go git make protobuf bash

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p /tmp/apate
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

RUN apk add curl --no-cache && \
    curl -L -o /usr/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl && \
    chmod +x /usr/bin/kubectl

RUN curl -L -o /helm.tar.gz https://get.helm.sh/helm-v3.2.1-linux-amd64.tar.gz && \
    tar -zxvf /helm.tar.gz && \
    mv /linux-amd64/helm /usr/bin/helm && \
    rm -rf /helm.tar.gz /linux-amd64

# JUnit report
RUN go get -v -u github.com/jstemmer/go-junit-report

# Controller-gen
RUN go get -v -u sigs.k8s.io/controller-tools/cmd/controller-gen

# Mockgen
RUN GO111MODULE=on go get -v -u github.com/golang/mock/mockgen@latest

# Proto
RUN GO111MODULE=on go get google.golang.org/grpc@v1.28.1
RUN GO111MODULE=on go get github.com/golang/protobuf/protoc-gen-go

WORKDIR $GOPATH
