# Implementation
This article will disucss the implementation decisions we have made and why we have made them.

## Virtual Kubelet
One of the biggest decisions to be made was how to emulate the Kubelets. It would have been possible to make our own Kubelet, but we chose to use Virtual Kubelet due to the fact that it is an existing product and hence is well-tested. Unfortunately, barebones Virtual Kubelet still requires a lot of boilerplate, given that the exact boilerplate we needed was hidden in their `internal` package. 

This is why we also used their `node-cli` project, even though it is not recommended to be used in production. We believe this is fine given that this system will only be used for testing purposes, and shouldn't serve any live applications.

## CRD
We use scenarios to define when, what and where different types of events should occur. These events can for example be node failure. In our initial implementation, we used our own format to define these and we made sure every node in the cluster had access to these scenarios so it knew what to do. The downside to this is that is non-standard and requires using non-standard tools.

After comparing different options, we came to the conclusion that CRDs would be best suited to replace these scenarios. A few upsides are that Kubernetes manages the validation, deserialization and distribution of deployments of these custom resources. Besides that, it is also easier for the end user, given that there are clear definitions of the structure in Kubernetes and you can thus use tools built for Kubernetes to control the emulation. 

A potential downside to this is that performance tests may not be completely fair, given that we also use the system under test. To remedy this, we recommend the usage of scenarios for performance tests, since they don't require updating the deployments in order to activate changes in the system. In fact, the deployments only need to be retrieved once before the performance test starts, and after that the Kubelets won't have to check for updated deployments anymore. This removes almost all overhead from CRDs with respect to our own solution, leading to fair tests.

## KinD
Setting up an entire Kubernetes cluster can be time consuming, so we provide a cluster to test with. This cluster uses KinD, which is basically a Kubernetes cluster inside of a Docker container. Even though this brings a bit of overhead, we think this is the easiest solution, only requiring Docker to be installed and thus not relying on a certain type of virtualization software. 

Because of the overhead, performance tests shouldn't be done using this provided cluster. In order to use a different cluster than KinD, one can provide the credentials for that cluster in the environment variables of the controlplane.

## gRPC & protobuf
The API endpoints of both the controlplane and the Apatelets use gRPC an protobuf. There are a few reasons to choose gRPC/protobuf above a REST/JSON API. First of all, gRPC provides two-way communication using streams. This is a big advantage for things like health services. Second, protobuf provides a well-defined request structure, something that JSON doesn't provide out-of-the-box. This and the fact that there are clients that can deal with protobuf for many languages made it a good choice for us.

The biggest downside of this approach is that the serialized version of protobuf messages are header to read and create than their JSON counterparts, making it harder to quickly test your API by for example using CURL to send a request.

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
