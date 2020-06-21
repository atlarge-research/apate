---
home: true
heroText: Apate
heroImage: /logo.png
tagline: "Low-resource Kubernetes emulation"
actionText: Get Started →
actionLink: /usage/
---

# Documentation index
We suggest you read the README of this project before reading these documents, as some terms are defined there which will be used without further introductions.

* [Structure](./structure.md)
* [Implementation choices](./implementation.md)

## Infrastructure
The application has 3 main components, which will be discussed here.

### Apatelet
An Apatelet is a "fake" Kubelet; it talks to Kubernetes as if it is a real Kubelet, but it doesn't end up actually running the pods, it will only tell Kubernetes it does. 
How this Apatelet and the pods running on the Apalet act can be controlled using a [CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).
This will be further discussed under [scenario](#scenario).

### Controlplane
This is the part which keeps track of which Apatelets are currently being run, and is able to spawn new Apatelets and kill old ones. You can choose to spawn these Apatelets using Docker containers or goroutines. Goroutines are faster and thus the default. It also keeps track of the health of the Apatelets and is able to determine when to stop the scenario. While the controlplane doesn't do much with the scenario, it is the middleman between the user and the Apatelets. The controlplane has a gRPC API through which it can be controlled.

### CLI
The CLI allows you to control the controlplane. It can create, get the status from and send updates to the controlplane. This just calls the gRPC API provided by the controlplane.

## Scenario
In order to control how the Apatelets and pods act towards Kubernetes, we use a "scenario". This scenario consists of different "tasks" which can be applied on the Apatelet / pods running on the Apatelet.

This task has a relative timestamp at which it will be run, with respect to the start of the scenario. This allows you to carefully define what happens to the cluster at a certain point in time.

Besides this scenario, you can also directly control the cluster, so state changes are applied immediately. This can also be done by applying a configuration based on the [CRDs we provide](configuration.md).

Examples can be found [here](examples.md).
