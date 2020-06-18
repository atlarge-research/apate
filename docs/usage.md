# Apate Usage Guide

[[toc]]

## Prerequisites
The following prerequisites are required for running Apate:
* Docker
* Modern Linux Distro

::: details OS Support  
We have only tested Apate on Linux. However it is likely other operating systems will work as long as Docker is present.
:::

::: warning  
Docker is technically not required as you could run the binaries directly, not using the Docker containers is heavily discouraged and not supported.
:::

## Installation
There are two ways of installing Apate. The recommended way is just by downloading and installing the pre-built binary. The other way is to download the source code and compile the binary.

### Pre-built binary <Badge text="recommended"/>
To install the `apate-cli` simply download the binary and put it in a `bin` directory.

To install it globally use this command
<!-- TODO: add url -->
```sh
curl -LO TODO_URL /usr/bin/apate-cli 
```

### From source
To install the `apate-cli` from source first clone the repo:

```bash
git clone https://github.com/atlarge-research/opendc-emulate-kubernetes
```

then move into the dir:
```bash
cd opendc-emulate-kubernetes
```

Finally you can install the `apate-cli` like this:
```bash
go install ./cmd/apate-cli/
```

Now `apate-cli` should be available under `$GOPATH/bin/apate-cli`.

::: tip  
It is also possible to build the controlplane and run it outside of a docker container. However that is out of scope for this guide.
:::

## Starting Apate
After you've install Apate and made sure the binary is present in your `$PATH` it is time to run Apate.

Starting up a kubernetes cluster and Apate control plane is as easy as just running:
```sh
apate-cli create
```

After this exits sucessfully you are ready to use Apate.

::: tip
It is also possible to use an already existing kubernetes cluser by providing `--kubeconfig-location` to the create command.
:::

### Kubeconfig
If you want to get the kubeconfig to connect to Apate run the following command:
```sh
apate-cli kubeconfig > ~/.kube/config
```
::: warning  
This _will_ overwrite your existing kube config, back it up if necessary
:::
<!-- TODO: Could we use contexts here to make this nicer to work with? -->

## Using Apate
Now that Apate is up and running it is time to run some experiments!

TODO: explain two emu modes

### Direct Emulation
Direct Emulation is the easiest way of experimenting with Apate and is done by applying and modifying various CRDs.
#### Nodes
First up we need to spin up some fake nodes (Apatelets). 

To do so first make a `.yml` with appropriate settings:
```yml
# ./example_node.yml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: test-deployment1
spec:
    replicas: 1
    resources:
        memory: 8G
        cpu: 4000
        storage: 3T
        ephemeral_storage: 500G
        max_pods: 110"
```

Now just use `kubectl` to spawn these Apatelets:

```sh
kubectl create -f ./example_node.yml
```

 To view all configuration options with explanations have look at the relevant sections of the [CRD Documentation](./configuration.md#nodes).


#### Pods
Now that we have spawned some nodes we are ready to run some pods.

You could just run some pods now and Apate would emulate them, however they wouldn't do anything and would perpetually in the running state.

To configure pod behaviour we need to create a custom PodConfiguration resource.

```yml
# ./example_pod.yml
apiVersion: apate.opendc.org/v1
kind: PodConfiguration
metadata:
    name: crd-deployment
spec:
	# TODO
```



### Planned Emulation
