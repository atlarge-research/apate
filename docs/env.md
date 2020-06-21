# Environment variables
Apate relies on environment variables for its configuration. All of them can be changed to fit to your needs. 

::: warning  
Some environment variables should almost never be changed, these are marked with a `*`.   
Only change these when you know what you are doing!   
:::  

[[toc]]

## Control plane
| Variable name | Type | Description | Default value |  
| --- | --- | --- | --- |   
| CP_LISTEN_ADDRESS | IP Address | Listen address for the control plane | 0.0.0.0 |   
| CP_LISTEN_PORT | Integer | Listen port for the control plane | 8085 |   
| CP_EXTERNAL_IP | IP Address | IP the Apatelets will attempt to connect to. Auto means the control plane will try to detect its own IP |  auto |   
| CP_DOCKER_POLICY | [Docker policy](#docker-policy) | Docker pull policy for apatelet images | pull-if-not-local |   
| CP_PROMETHEUS | Boolean | Enable prometheus installation | true |   
| CP_PROMETHEUS_NAMESPACE | String | The namespace where prometheus will be installed | apate-prometheus |   
| CP_KUBE_CONFIG | String | The kube config used for interacting with the Kubernetes cluster | \<unset> |   
| CP_KUBE_CONFIG_LOCATION | String | Path to the kube config | /tmp/apate/config |   
| CP_APATELET_RUN_TYPE | [Run type](#run-type) | The run type used for running new Apatelets | ROUTINES |   
| CP_KIND_CLUSTER_NAME | String | Cluster name | apate |   
| CP_PROMETHEUS_CONFIG_LOCATION*  | String | The path to the config of the prometheus helm chart, if applicable | config/prometheus.yml |   
| CP_POD_CRD_LOCATION* | String | Location of the pod CRD | config/crd/apate.opendc.org_nodeconfigurations.yaml |   
| CP_NODE_CRD_LOCATION* | String | Location of the node CRD | config/crd/apate.opendc.org_podconfigurations.yaml |   
| CP_DOCKER_HOSTNAME* | String | Enable hostname rewriting to 'docker' | false |   
| CP_ENABLE_DEBUG* | Boolean | Enable extra debug messages | false |   
| CP_MANAGER_CONFIG_LOCATION* | String | Path to the config of the cluster manager, if applicable | config/kind.yml |   

### Types
#### Docker policy
When docker containers are being used, these types define how the control plane handles missing images.   

| Name | Description | 
| --- | --- |
| pull-always | AlwaysPull will always pull the image, even if it's already locally available |
| local-always | AlwaysLocal will always use the local image, and will not pull if it's not locally available |
| pull-if-not-local | PullIfNotLocal will only pull if the image is not locally available, this will not check if the local image is possibly outdated |

#### Run type
Apatelets can be spawned using different runners. Depending on the type of runner, the Apatelet might be spawned in a container, or just an extra process.   

| Name | Description | 
| --- | --- |
| ROUTINES | This will run the Apatelets in seperate go routines. This has less overhead and will probably faster, however, this will greatly reduce the usefullness of logs. |
| DOCKER | This will run each Apatelet in its own Docker container. This will limit the amount of Apatelets to about 600 instances on most machines, due to Docker restrictions. You could work around this, but for large tests we suggest using the routine runner. |

## Apatelet
Most information the Apatelet needs it transfered to it, once it joins the Apate cluster. However, some key pieces of information
are needed during initialization, or are more easily calculated in the runner. 
These environment variables are set by the Apatelet runner upon creation of a new Apatelet. If you plan on creating a new runner,
these variables are important!

::: danger  
All environment variables are set by the Apatelet runner when spawning Apatelets.
Changing these outside of the runner might have unforeseen consequences!  
:::  

| Variable name | Type | Description | Default value |  
| --- | --- | --- | --- |   
| APATELET_LISTEN_ADDRESS | IP Address | Listen address for the Apatelet | 0.0.0.0 |
| APATELET_LISTEN_PORT | Integer | Listen port for the Apatelet | 8086 |
| APATELET_METRICS_PORT | Integer | Listen port for Kubernetes daemon | 10250 |
| APATELET_KUBERNETES_PORT | Integer | Listen port for Prometheus metrics | 10255 |
| APATELET_KUBE_CONFIG | String | Path to the kube config | /tmp/apate/config |
| APATELET_CP_ADDRESS | IP Address | Address that should be used when connecting to the control plane  | localhost |
| APATELET_CP_PORT | Integer | Port that should be used when connecting to the control plane | 8085 |
| APATELET_DISABLE_TAINTS | Boolean | Determines whether to disable taints on this node or not | false |
| APATELET_ENABLE_DEBUG* | Boolean | Enable extra debug messages | false |
| CI_APATELET_K8S_ADDRESS* | String | CIKubernetesAddress is to add an entry to the /etc/hosts file of an apatelet to ensure it can find the k8s cluster. This can be used to fix bugs when running apate in a DinD | \<unset> |