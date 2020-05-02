GO111MODULE=on

all: build protobuf mockgen lint test docker_build

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
	go test -v ./...

test_short:
	go test -short ./...

test_race:
	go test -race -short ./...

test_cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

test_cover_short:
	go test -short -coverprofile cover.out ./...
	go tool cover -html=cover.out

docker_build_vk:
	docker build -f services/apatelet/Dockerfile -t apatelet .

docker_build_cp:
	docker build -f ./services/controlplane/Dockerfile -t controlplane .

docker_build: docker_build_cp docker_build_vk

run_cp: docker_build_cp
	docker run --network host -v /var/run/docker.sock:/var/run/docker.sock controlplane

protobuf:
	protoc -I ./api --go_opt=paths=source_relative --go_out=plugins=grpc:./api/ `find . -type f -name "*.proto" -print`

# Generates the various mocks
mockgen: ./api/health/mock_health/health_mock.go ./services/controlplane/store/mock_store/store_mock.go ./services/apatelet/store/mock_store/store_mock.go

./api/health/mock_health/health_mock.go: ./api/health/health.pb.go
	mockgen github.com/atlarge-research/opendc-emulate-kubernetes/api/health Health_HealthStreamClient,HealthClient,Health_HealthStreamServer > $@

./services/controlplane/store/mock_store/store_mock.go: ./services/controlplane/store/store.go
	mockgen github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store Store > $@

./services/apatelet/store/mock_store/store_mock.go: ./services/apatelet/store/store.go
	mockgen github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store Store > $@
