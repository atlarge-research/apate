# Usage Guide

[[toc]]

## Prerequisites
The following prerequisites are required for running Apate:
* Docker
* Modern Linux Distro

::: details OS Support  
We have only tested Apate on Linux. However it is likely other operating systems will work as long as Docker is present.
:::

## Installation
There are two ways of installing Apate. The recommended way is just by downloading and installing the pre-built binary. The other way is to download the source code and compile the binary.

### Pre-built binary <Badge text="recommended"/>
To install the `apate-cli` globally simply download the binary and put it in a `bin` directory:
```sh
curl -LO https://github.com/atlarge-research/opendc-emulate-kubernetes/releases/download/v0.1.0/apate-cli /usr/bin/apate-cli
```

You can also use the CLI relatively:
```sh
curl -LO https://github.com/atlarge-research/opendc-emulate-kubernetes/releases/download/v0.1.0/apate-cli
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

::: tip Note
It is also possible to build the controlplane and run it outside of a Docker container. However that is out of scope for this guide.
:::

## Starting Apate
Starting up a kubernetes cluster and Apate control plane is as easy as just running:
```sh
apate-cli create # or ./apate-cli create if used relatively
```

After this exits sucessfully you are ready to use Apate.

::: tip
It is also possible to use an already existing kubernetes cluster by providing `--kubeconfig-location` to the create command.
:::

::: details Starting Control plane manually
You can start the control plane Docker container manually like so: `docker run -v /var/run/docker.sock:/var/run/docker.sock -p 8085:8085 apatekubernetes/controlplane`
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

### Nodes
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
        max_pods: 110
```

Now just use `kubectl` to spawn these Apatelets:

```sh
kubectl create -f ./example_node.yml
```

Now you can use `kubectl get nodes` to see if the node was actually created. It should look something like this:
```
NAME                                            STATUS   ROLES    AGE     VERSION
apate-control-plane                             Ready    master   6m51s   v1.17.0
apatelet-f97ae5b6-97d3-4e6f-8f20-bfbe724a8bd2   Ready    agent    2s      v1.15.2
```

::: tip  
To view all configuration options with explanations have look at the relevant sections of the [CRD Documentation](./configuration.md#nodes).
:::

### Pods
Now that we have spawned some nodes we are ready to run some pods.

You could just run some pods now and Apate would emulate them, however they wouldn't do anything and would perpetually in the running state.

To configure pod behaviour we need to create a custom [`PodConfiguration`](./configuration.md#pods) resource.

```yml
# ./example_pod.yml
apiVersion: apate.opendc.org/v1
kind: PodConfiguration
metadata:
    name: crd-deployment
spec:
    pod_status: "SUCCEEDED"
```

This status dictates that a pod should succeed immediately.

Now we need to make a config for some pods actually using this CRD.

This configuration is fairly typical except for a couple of points.
```yml
# ./emulated_nginx.yml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        apate: crd-deployment
    spec:
      nodeSelector:
        emulated: "yes"
      tolerations:
        -   key: emulated
            operator: Exists
            effect: NoSchedule
      containers:
        - name: nginx
          image: nginx:1.14.2
          ports:
            - containerPort: 80
```
The `apate: crd-deployment` label is to select which CRD this pod should use and should therefore match the name you defined in the PodConfigration.
The other difference from your typical deployment is the `nodeSelector` and `tolerations` these are to specify that this pod should be ran on emulated nodes.

Now let's run this:
```sh
kubectl apply -f example_pod.yml
kubectl apply -f emulated_nginx.yml
```

If you now check the status of these pods you should see the pods have succeeded.
```
NAME                                READY   STATUS      RESTARTS   AGE
nginx-deployment-6cbf99798f-55b9k   0/1     Succeeded   0          37s
nginx-deployment-6cbf99798f-5jrdr   0/1     Succeeded   0          18s
nginx-deployment-6cbf99798f-62ftl   0/1     Pending     0          3s
```

Because this is a deployment you can see the controller is trying to start new pods up covering for the succeeded ones. Make sure to stop the deployment using `kubectl delete -f emulated_nginx.yml` otherwise you might end up with a lot of pods.

::: tip Emulation Methods  
This example used so-called direct emulation, to view how to do planned emulation and for other configuration options please look at the [CRD Configuration](./configuration.md) page.
:::