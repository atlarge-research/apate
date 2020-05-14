GO111MODULE=on

all: build protobuf mockgen lint test docker_build

.PHONY: build
build:
	go build ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint_fix_imports
lint_fix_imports:
	goimports -local github.com/atlarge-research/opendc-emulate-kubernetes -w **/*.go
	golangci-lint run

.PHONY: lint_fix
lint_fix:
	golangci-lint run --fix

.PHONY: test
test:
	go test -v ./...

.PHONY: test_short
test_short:
	go test -short ./...

.PHONY: test_race
test_race:
	go test -race -short ./...

.PHONY: test_cover
test_cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

.PHONY: test_cover_short
test_cover_short:
	go test -short -coverprofile cover.out ./...
	go tool cover -html=cover.out

.PHONY: docker_build_vk
docker_build_vk:
	docker build -f services/apatelet/Dockerfile -t apatelet .
	docker tag apatelet apatekubernetes/apatelet

.PHONY: docker_build_cp
docker_build_cp:
	docker build -f ./services/controlplane/Dockerfile -t controlplane .
	docker tag controlplane apatekubernetes/controlplane

docker_build: docker_build_cp docker_build_vk

.PHONY: docker_build_cp
run_cp: docker_build_cp
	docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8085:8085 controlplane

.PHONY: protobuf
protobuf:
	protoc -I ./api --go_opt=paths=source_relative --go_out=plugins=grpc:./api/ `find . -type f -name "*.proto" -print`

# Generates the various mocks
mock_gen: ./api/health/mock_health/health_mock.go ./services/controlplane/store/mock_store/store_mock.go ./services/apatelet/store/mock_store/store_mock.go ./services/apatelet/provider/mock_cache_store/mock_cache_store.go

./api/health/mock_health/health_mock.go: ./api/health/health.pb.go
	mockgen github.com/atlarge-research/opendc-emulate-kubernetes/api/health Health_HealthStreamClient,HealthClient,Health_HealthStreamServer > $@

./services/controlplane/store/mock_store/store_mock.go: ./services/controlplane/store/store.go
	mockgen github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store Store > $@

./services/apatelet/store/mock_store/store_mock.go: ./services/apatelet/store/store.go
	mockgen github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store Store > $@

./services/apatelet/provider/mock_cache_store/mock_cache_store.go: FORCE
	mockgen k8s.io/client-go/tools/cache Store > $@

crd_gen:
	controller-gen object paths=./pkg/apis/emulatedpod/...
	controller-gen crd:trivialVersions=false,crdVersions=v1 paths=./pkg/apis/emulatedpod/...

gen: crd_gen mock_gen protobuf

FORCE:
