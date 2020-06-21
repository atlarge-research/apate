# Build Guide
## Regenerating code
To regenerate generated code, which you don't have to do unless you protobuf / structs which are mocked, you need a few dependencies, namely: 
* `controller-gen`, which can be installed by running `go get -v -u sigs.k8s.io/controller-tools/cmd/controller-gen`
* [`mockgen`](https://github.com/golang/mock)
* [`protobuf`](https://grpc.io/docs/languages/go/quickstart/)

## Building
It is easy to build the project, using the commands defined in the Makefile. We have made a Dockerfile which bundles the application and all its dependencies, which can be called using `make docker_build`.