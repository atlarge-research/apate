GO111MODULE=on

all: build lint test

build:
	go build ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./...

test-race:
	go test -race -short ./...
