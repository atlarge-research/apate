GO111MODULE=on

all: build lint test

build:
	go build ./...

lint:
	golangci-lint run

lint_fix_imports:
	goimports -local github.com/atlarge-research/opendc-emulate-kubernetes -w **/*.go
	golangci-lint run

lint_fix:
	golangci-lint run --fix

test:
	go test ./...

test_short:
	go test -short ./...

test_race:
	go test -race -short ./...

docker_build_vk:
	docker build -f services/apatelet/Dockerfile -t apatelet .

docker_build_cp:
	docker build -f ./services/controlplane/Dockerfile -t controlplane .

docker_build: docker_build_cp docker_build_vk

run_cp: docker_build_cp
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock controlplane
