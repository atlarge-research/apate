# CRD Configuration
Apate makes uses of [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) to provide an easy way of configuring the emulated nodes and pods. 

## Planned and direct emulation
Apate supports both planned and direct emulation. Both have their advantages and disadvantages, and serve specific use-cases.

Planned emulation is our regular scenario. It consists of a list of tasks that should be executed at a specific time relative 
to the start of the scenario. 

::: warning
Even though most use-cases are covered by either planned or direct emulation, some edge-cases might require both.
However, this is not supported in the initial version of Apate, and might lead to unexpected results in some cases.  
We invite others to contribute to Apate and add this feature, as it should be a good first issue.
:::  

## Nodes
A `NodeConfiguration` describes a set of emulated nodes in the Kubernetes cluster. The specification describes a list of tasks 
for the scenario, a resource definition, and the amount of nodes. Optionally, one can provide a direct state, as opposed to a list
of tasks. 

For example, the following `NodeConfiguration` will create ten worker nodes, all with 5G memory, 5 cores, 5T storage, 120G ephemeral storage,
and a maximum of 150 pods. All nodes will 'fail' one second after the scenario has started:

```yaml
apiVersion: apate.opendc.org/v1
kind: NodeConfiguration
metadata:
    name: test-deployment
spec:
    replicas: 10
    resources:
        memory: 5G
        cpu: 5
        storage: 5T
        ephemeral_storage: 120G
        max_pods: 150
    tasks:
        - timestamp: 1s
          state:
              node_failed: true
```

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| \<inline> | [State](#node-state) | A state to immediately apply| No |
| replicas | int64 | The amount of nodes that are required| Yes
| resources | [Resources](#node-resources) | A resource specification| Yes |
| tasks | [Task\[\]](#node-task) | A list of tasks for these nodes | No |

### Node resources
Resources describe the amount of emulated resources this node has.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| memory | [Bytes](#bytes) | Amount of memory | Yes |
| cpu | int64 | Amount of [CPU cores](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#meaning-of-cpu) | Yes |
| storage | [Bytes](#bytes) | Amount of storage | Yes |
| ephemeral_storage | [Bytes](#bytes) | Amount of ephemeral storage | Yes | 
| max_pods | in64 | Maximum amount of pods | Yes |


### Node task
Task is a combination of a timestamp and a state.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| timestamp | [Time](#time) | Time at which this task will be executed| Yes |
| state | [State](#node-state) | Desired state after this task | Yes |

### Node state
State is the desired state of the node. 

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| node_failed | bool | If true, will no longer react to any requests and will no longer send heartbeats | No |
| network_latency | [Time](#time) | Applies extra latency to request from Kubernetes| No |
| heartbeat_failed | bool | If true, will no longer send heartbeats to Kubernetes| No |
| custom_state | [Custom state](#custom-state) | A custom state | No |

::: warning  
In the initial version of Apate, it is not possible to revert `node_failed` or `heartbeat_failed`. 
We invite others to contribute to Apate and add this feature, as it should be a good first issue.  
:::

#### Custom state
Custom state can be used to directly modify the internal flags. This allows users to create a custom state.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| create_pod_response | [Response](#response) | Response to creation of pods | No |
| update_pod_response | [Response](#response) | Response to updates of pods | No | 
| delete_pod_response | [Response](#response) | Response to deletion of pods  | No |
| get_pod_response | [Response](#response) | Response to pod requests  | No |
| get_pod_status_response | [Response](#response) | Response to pod status requests | No |
| node_ping_response | [Response](#response) | Response to node heartbeats | No |

## Pods
A `PodConfiguration` describes a set of emulated pods in the Kuber  netes cluster. The specification simply contains a list 
of tasks for the pods. Optionally, one can provide a direct state, as opposed to a list of tasks. 

For example, the following `PodConfiguration` will ensure all pods with the label pair `<key, value>` (given this CRD is
deployed in the namespace `namespace`, start with a memory usage of 1G. After one second has passed after the start of the 
scenario, its usage will increase to 5G of memory and 1 core. Finally, it will fail after five seconds:
```yaml
apiVersion: apate.opendc.org/v1
kind: PodConfiguration
metadata:
    name: crd-deployment
spec:
    pod_resources:
        memory: 1G
    tasks:
        - timestamp: 1s
          state:
              pod_resources:
                  cpu: 1
                  memory: 5G
        - timestamp: 5s
          state:
              pod_status: FAILED
```

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| \<inline> | [State](#pod-state) | A state to immediately apply| No |
| tasks | [Task\[\]](#pod-task) | A list of tasks for these pods | No |

### Pod resources
Resources describe the amount of emulated resources this pod uses.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| memory | [Bytes](#bytes) | Amount of memory | No |
| cpu | int64 | Amount of CPU | No |
| storage | [Bytes](#bytes) | Amount of storage | No |
| ephemeral_storage | [Bytes](#bytes) | Amount of ephemeral storage | No | 

### Pod task
Task is a combination of a timestamp and a state

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| timestamp | [Time](#time) | Time at which this task will be executed | Yes |
| relative_to_pod | bool | If true, the timestamp will be relative to the start time of the pod, instead of the start time of the scenarion | No |
| state | [State](#pod-state) | Desired state after this task | Yes |

### Pod state
State is the desired state of the pod.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| create_pod_response | [Response](#response) | Response to creation of pods | No |
| update_pod_response | [Response](#response) | Response to updates of pods | No | 
| delete_pod_response | [Response](#response) | Response to deletion of pods  | No |
| get_pod_response | [Response](#response) | Response to pod requests  | No |
| get_pod_status_response | [Response](#response) | Response to pod status requests | No |
| pod_resources | [Resources](#pod-resources) | Pod resource usage | No |
| pod_status | [Status](#status) | Pod status | No |

## Types
To more easily work with our CRD, we have added a few extra types.

### Time
This type can be used to easily work with time. This should be appended to an integer.

| Type | Description |
| --- | --- |
| ns | Nanosecond |
| us | Microsecond |
| ms | Millisecond |
| s | Second |
| m | Minute |
| h | Hour |

### Bytes
This type can be used to easily work with bytes. This should be appended to an integer.

| Type | Description |
| --- | --- |
| k | KiB / Kibibyte |
| m | MiB / Mebibyte |
| g | GiB / Gibibyte |
| t | TiB / Tebibyte |
| p | PiB / Pebibyte |

### Response
This type can be used to easily work with response types.

| Type | Description |
| --- | --- |
| NORMAL | The request will be answered normally |
| TIMEOUT | The request will timeout |
| ERROR | The node will return an error | 

### Status
This type can be used to determine the status of a pod.

| Type | Description |
| --- | --- |
| PENDING | Pod is pending |
| RUNNING | Pod is running |
| SUCCEEDED | Pod stopped successfully |
| FAILED | Pod stopped with an exit code signaling failure |  
| UNKNOWN | Pod status is unknown |