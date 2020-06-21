# Emulating Kubernetes
Apate deals with emulating Kubernetes pods, making it possible to run thousands of pods on a single laptop to emulate how Kubernetes will respond to that.

## Documentation
The full Apate documentation can be found [here](https://apatekubernetes.nl/)

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
