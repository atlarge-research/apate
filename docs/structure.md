# Structure
This article will discuss the structure and workflow of Apate and explain why we made the choices to use such a workflow.

## Workflow
![Apate's architecture](../documents/Architecture_flow.svg)

### Starting the control plane
First, the CLI can be used to start the controlplane. Now the CLI will be waiting and polling the controlplane and ask for the amount of Apatelets that are healthy. The controlplane store will contain references to the Apatelets that have been started and their health state (e.g. healthy).

### Spawning Apatelets
No Apatelets will are spawned until the user asks for them by creating a `NodeConfiguration`. When the user asks for new Apatelets by for example using `kubectl create -f node-configuration.yml`, Kubernetes will tell the informer running on the controlplane that this CRD has been applied. The controlplane will spawn Apatelets accordingly, either on Docker or as a goroutine.

### Apatelet start
When an Apatelet is spawned, it will ask the controlplane for the kubeconfig, which is required to join the Kubernetes cluster. It will also get its hardware specification, which will be sent to Kubernetes. With this information, it tries to join the Kubernetes cluster. When this is succesful, it will start its health service and it will tell the controlplane that it's healthy. Now the controlplane will update its store to mark the Apatelet as healthy (which will update the response from the status service).

### Spawning pods & their configurations
Just like nodes, there is a CRD to configure pods. For this purpose, every Apatelet listens to updates to the `PodConfiguration` CRD. When it is informed by such an update, it will update its internal state accordingly. More on this [here](./implementation.md).

### Starting the scenario
After configuring the nodes and pods, we can actually start the scenario. This can be done by using the CLI, which will call the controlplane, which will in turn call every Apatelet to tell it to start the scenario.

## CRD
We use scenarios to define when, what and where different types of events should occur. These events can for example be node failure. In our initial implementation, we used our own format to define these and we made sure every node in the cluster had access to these scenarios so it knew what to do. The downside to this is that is non-standard and requires using non-standard tools.

After comparing different options, we came to the conclusion that CRDs would be best suited to replace these scenarios. A few upsides are that Kubernetes manages the validation, deserialization and distribution of deployments of these custom resources. Besides that, it is also easier for the end user, given that there are clear definitions of the structure in Kubernetes and you can thus use tools built for Kubernetes to control the emulation. 

A potential downside to this is that performance tests may not be completely fair, given that we also use the system under test. To remedy this, we recommend the usage of scenarios for performance tests, since they don't require updating the deployments in order to activate changes in the system. In fact, the deployments only need to be retrieved once before the performance test starts, and after that the Kubelets won't have to check for updated deployments anymore. This removes almost all overhead from CRDs with respect to our own solution, leading to fair tests.
