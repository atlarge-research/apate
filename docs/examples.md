# Examples
For your convenience we provide some small examples here. This should be able to help you understand Apate's configuration.
If not, feel free to create an issue and we are glad to help!

[[toc]]

## Simple cluster
In this example, we will spin up a simple kubernetes cluster of ten worker nodes. There will be no node failure or anything of the sort.
This cluster can be used to simply emulated pods.

```yaml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: simple-node-deployment
spec:
    replicas: 10
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
```

## Massive node failure
This example will simulate a cluster of 100 nodes. Five seconds after the scenario has started, 80 of those nodes will fail.

```yaml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: massive-node-failure-1
spec:
    replicas: 80
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
    tasks:
        - timestamp: 5s
          state:
              node_failed: true
---
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: massive-node-failure-2
spec:
    replicas: 20
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
```

## Pod failure
In this example, we will create a simple cluster of three worker nodes. We will run an nginx deployment on this cluster, which 
will fail ten seconds after the scenario has started.  
Note that Kubernetes will try to reschedule the failed pods, but they will fail again, resulting in an infinite loop. 
This is expected behaviour for deployments.

```yaml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: pod-failure
spec:
    replicas: 3
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
---
apiVersion: apate.opendc.org/v1
kind: PodConfiguration
metadata:
    name: pod-failure
spec:
    tasks:
        - timestamp: 10s
          state:
              pod_status: FAILED
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-pod-failure
  labels:
    app: nginx
spec:
  replicas: 10
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
        apate: pod-failure
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

## Node failure with direct emulation
This example will showcase direct emulation. Initially we will create a simple deployment of one worker node with no tasks.
```yaml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: test-deployment-1
spec:
    replicas: 1
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
```

Once you are ready for node failure, you can simply change the file to include the direct state, and update the CRD in Kubernetes.
This will result in the node failing, without using the scenario.

```yaml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: test-deployment-1
spec:
    replicas: 1
    node_failed: true
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
```

## Slow network and node failure
For the final example, we will emulate a cluster of five worker nodes. Immediately after the scenario starts all of them
will suddenly have a slow network with an added delay of 100ms. Five seconds after this, one of the nodes will fail.

```yaml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: slow-network-1
spec:
    replicas: 1
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
    tasks:
        - timestamp: 0s
          state:
              network_latency: 100ms
        - timestamp: 5s
          state:
              node_failed: true
---
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: slow-network-2
spec:
    replicas: 4
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
    tasks:
        - timestamp: 0s
          state:
              network_latency: 100ms
```