# Implementation
This article will disucss the implementation decisions we have made and why we have made them.

## Virtual Kubelet
One of the biggest decisions to be made was how to emulate the Kubelets. It would have been possible to make our own Kubelet, but we chose to use Virtual Kubelet due to the fact that it is an existing product and hence is well-tested. Unfortunately, barebones Virtual Kubelet still requires a lot of boilerplate, given that the exact boilerplate we needed was hidden in their `internal` package. 

This is why we also used their `node-cli` project, even though it is not recommended to be used in production. We believe this is fine given that this system will only be used for testing purposes, and shouldn't serve any live applications.

## KinD
Setting up an entire Kubernetes cluster can be time consuming, so we provide a cluster to test with. This cluster uses KinD, which is basically a Kubernetes cluster inside of a Docker container. Even though this brings a bit of overhead, we think this is the easiest solution, only requiring Docker to be installed and thus not relying on a certain type of virtualization software. 

Because of the overhead, performance tests shouldn't be done using this provided cluster. In order to use a different cluster than KinD, one can provide the credentials for that cluster in the environment variables of the controlplane.

## gRPC & protobuf
The API endpoints of both the controlplane and the Apatelets use gRPC an protobuf. There are a few reasons to choose gRPC/protobuf above a REST/JSON API. First of all, gRPC provides two-way communication using streams. This is a big advantage for things like health services. Second, protobuf provides a well-defined request structure, something that JSON doesn't provide out-of-the-box. This and the fact that there are clients that can deal with protobuf for many languages made it a good choice for us.

The biggest downside of this approach is that the serialized version of protobuf messages are header to read and create than their JSON counterparts, making it harder to quickly test your API by for example using CURL to send a request.

## Apatelet store
Each Apatelet has a store which holds a number of flags. Each flag influences how the Apatelet behaves.

For example: one flag might be to error on all `GetPod` requests sent by Kubernetes, or to add a certain amount of latency. We think flags are a good option here because you can combine multiple flags to form "tasks". An example task is node failure, which just times out on all requests sent to it. Doing this allows us to emulate real world scenarios such as node failure or network segregation, and allows us to see how Kubernetes reacts to this.

These stores are fully thread safe, making sure that simultaneous updates to the store are handled properly. In order to actually update the store, we use a scheduler, which "executes" the tasks which have been added to it by the CRDs. You can also directly update these flags in order to update how the Apatelet responds to requests.

## Extendability
In order to make it possible to easily change a few aspects of our application we have tried to abstract away a lot of functionality. There are two most notable locations in which we have done this.

#### Spawning Apatelets
We wanted to provide multiple ways to spawn Apatelets and make it easy for an external user to add their own implementation. In order to do that, we have made an interface which needs to be implemented in order to spawn Apatelets. We have made implementations for spawning them as goroutines and spawning them as Docker containers on the host, but it is easy to provide your own implementation. // TODO add name of interface

#### Cluster
As explained previously, we provide the possibility to create a new Kubernetes cluster when starting up the control plane, but we also provide other ways to connect to a cluster. For that, you only need to implement an interface. // TODO add name of interface

## Delivery format
### CLI
We deliver the CLI as an executable. The only dependency it has is that Docker is installed on the system.

### Controlplane
We only deliver the controlplane as a Docker container, as it has many dependencies.
