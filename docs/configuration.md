# CRD Configuration
Apate makes uses of [CRDs](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) to provide an easy way of configuring the emulated nodes and pods. 

## Nodes
A `NodeConfiguration` describes a set of emulated nodes in the Kubernetes cluster. The specification describes a list of tasks 
for the scenario, a resource definition, and the amount of nodes. Optionally, one can provide a direct state, as opposed to a set
of tasks. 

For example, the following NodeConfiguration will create ten worker nodes, all with 5G of memory etc etc etc. 
All of the nodes will 'fail' one second after the scenario has started:

```yaml
todo: yes
```

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| \<inline> | [State](#state) | A state to immediately apply| No |
| replicas | int64 | The amount of nodes that are required| Yes
| resources | [Resources](#resources) | A resource specification| Yes |
| tasks | [Task\[\]](#task) | A list of tasks for these nodes | No |

### Resources
Resources describes the amount of emulated resources this node has.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| memory | [Bytes](#bytes) | Amount of memory | Yes |
| cpu | int64 | Amount of CPU | Yes |
| storage | [Bytes](#bytes) | Amount of storage | Yes |
| ephemeral_storage | [Bytes](#bytes) | Amount of ephemeral storage | Yes | 
| max_pods | in64 | Maximum amount of pods | Yes |


### Task
Task is a combination of a timestamp and a state.

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| timestamp | [Time](#time) | Time at which this task will be executed| Yes |
| state | [State](#state) | Desired state after this task | Yes |

### State
State is the desired state of the node. 

| Field | Type | Description | Required |
| --- | --- | --- | --- |
| node_failed | bool | If true, will no longer react to any requests and will no longer send heartbeats | No |
| network_latency | int64 | Applies extra latency to request from Kubernetes| No |
| heartbeat_failed | bool | If true, will no longer send heartbeats to Kubernetes| No |
| custom_state | [Custom state](#custom-state) | A custom state | No |

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

## Pods