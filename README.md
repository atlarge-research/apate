# Emulating Kubernetes
This project deals with emulating Kubernetes pods, making it possible to run thousands of pods on a single laptop to emulate how Kubernetes will respond to that.

## Infrastructure
The application has 3 main components, which will be discussed here.

### Apatelet
The Apatelet is a "fake" Kubelet; it talks to Kubernetes as if it is a real Kubelet, but it doesn't end up actually running the pods, it will only tell Kubernetes it does. 
How this Apatelet and the pods running on the Apalet act can be controlled using a [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
This will be further discussed under [scenario](#scenario).

### Controlplane
This is the part which keeps track of which Apatelets are currently being run, and is able to spawn new Apatelets and kill old ones. You can choose to spawn these Apatelets using Docker containers or goroutines. Goroutines are faster and thus the default. It also keeps track of the health of the Apatelets and is able to determine when to stop the scenario. While the controlplane doesn't do much with the scenario, it is the middleman between the user and the Apatelets. The controlplane has a gRPC API through which it can be controlled.

### CLI
The CLI allows you to control the controlplane. It can create, get the status from and send updates to the controlplane. This just calls the gRPC API provided by the controlplane.

## Scenario
In order to control how the Apatelets and pods act towards Kubernetes, we use a "scenario". This scenario consists of different "tasks" which can be applied on the Apatelet / pods running on the Apatelet.

This task has a relative timestamp at which it will be run, with respect to the start of the scenario. This allows you to carefully define what happens to the cluster at a certain point in time.

Besides this scenario, you can also directly control the cluster, so state changes are applied immediately. This can also be done by applying a configuration based on the [CRDs we provide](config/crd).

Example CRDs can be found [here](examples).

## Build process
### Regenerating code
To regenerate generated code, which you don't have to do unless you protobuf / structs which are mocked, you need a few dependencies, namely: `controller-gen`, `mockgen` and `protobuf`.

### Building
It is easy to build the project, using the commands defined in the Makefile. We have made a Dockerfile which bundles the application and all its dependencies, which can be called using `make docker_build`.

## Installation
### CLI
In order to easily control the entire application, you can download the CLI from [Github](https://github.com/atlarge-research/opendc-emulate-kubernetes/releases).

### Docker
You can also download the Apate controlplane or the Apatelet from [dockerhub](https://hub.docker.com/u/apatekubernetes):

```bash
docker pull apatekubernetes/controlplane:latest
docker pull apatekubernetes/apatelet:latest
```

## Usage
### CLI
The CLI contains all the tools to start and manage the lifecycle of your application.

First off, you can create a local control plane by running `apate-cli create`. Printing out the options can be done by running `apate-cli create --help`

When this control plane is running, you can obtain the kubeconfig by running `apate-cli kubeconfig`.

### Docker
Running the controlplane can also be done manually using Docker. For this, you first pull the controlplane as described previously.

Then you can execute the Docker container with default variables using `docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8085:8085 controlplane`. The environment variables that can be provided can be found [here](pkg/env/controlplane.go).

## Documentation
Apate's documentation can be found [here](docs/index.md)

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

### Project structure
    .
    ├── api/          # gRPC definitions for services both on the Apatelet and the controlplane
    ├── ci/           # CI configuration files
    ├── cmd/          # The code for the CLI
    ├── config/       # Configuration files used by the control plane
    ├── docs/         # Documentation
    ├── examples/     # Some example CRDs
    ├── internal/     # Internal utilities 
    ├── pkg/          # Utilities that could be useful to anyone
    ├── services/     # Contains the Apatelet and controlplane main files
    ├── Makefile      # A makefile containing some commands which aid in building / testing the project
    ├── go.*          # The go module and sum files, containing versions of dependencies
    └── LICENSE       # The license of this project

## License
Apate is licensed under the [Apache 2.0 license](./LICENSE)
