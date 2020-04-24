GO111MODULE=on

all: build lint test

build:
	go build ./...

lint:
	golangci-lint run

lint-fix-imports:
	goimports -local github.com/atlarge-research/opendc-emulate-kubernetes -w **/*.go
	golangci-lint run

lint-fix:
	golangci-lint run --fix

test:
	go test ./...

test_short:
	go test -short ./...

test-race:
	go test -race -short ./...

docker_build:
	docker build -f services/virtual_kubelet/Dockerfile -t virtual_kubelet .
