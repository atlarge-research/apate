GO111MODULE=on

all: build lint test

build:
	go build ./...

lint:
	golangci-lint run

lint_fix:
	golangci-lint run --fix

test:
	go test ./...

test_race:
	go test -race -short ./...

docker_build_vk:
	docker build -f services/virtual_kubelet/Dockerfile -t virtual_kubelet .

docker_build_cp:
	docker build -f ./services/control_plane/Dockerfile -t control_plane .

docker_build: docker_build_cp docker_build_vk

run_cp: docker_build_cp
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock control_plane
