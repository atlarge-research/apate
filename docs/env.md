# Environment variables
Apate relies on environment variables for its configuration. All of them can be changed to fit to your needs. 

::: warning  
Some environment variables should almost never be changed, these will be <bold/italic/?>.   
Only change these when you know what you are doing!   
:::  

## Control plane
| Variable name | Description | Default value |  
| --- | --- | :---: |   
| CP_LISTEN_ADDRESS | Listen address for the control plane | 0.0.0.0 |   
| CP_LISTEN_PORT | Listen port for the control plane | 8085 |   
| CP_EXTERNAL_IP | IP the Apatelets will attempt to connect to. Auto means the control plane will try to detect its own IP |  auto |   
| CP_DOCKER_POLICY | Docker pull policy for apatelet images | pull-if-not-local |   
| CP_PROMETHEUS | Enable prometheus installation | true |   
| CP_PROMETHEUS_NAMESPACE | The namespace where prometheus will be installed | apate-prometheus |   
| CP_KUBE_CONFIG | The kube config used for interacting with the Kubernetes cluster |  |   
| CP_KUBE_CONFIG_LOCATION | Path to the kube config | /tmp/apate/config |   
| CP_APATELET_RUN_TYPE | The run type used for running new Apatelets | ROUTINES |   
| CP_KIND_CLUSTER_NAME | Cluster name | apate |   
| **CP_PROMETHEUS_CONFIG_LOCATION** | The path to the config of the prometheus helm chart, if applicable | config/prometheus.yml |   
| **CP_POD_CRD_LOCATION** | Location of the pod CRD | config/crd/apate.opendc.org_nodeconfigurations.yaml |   
| **CP_NODE_CRD_LOCATION** | Location of the node CRD | config/crd/apate.opendc.org_podconfigurations.yaml |   
| **CP_DOCKER_HOSTNAME** | Enable hostname rewriting to 'docker' | false |   
| **CP_ENABLE_DEBUG** | Enable extra debug messages | false |   
| **CP_MANAGER_CONFIG_LOCATION** | Path to the config of the cluster manager, if applicable | config/kind.yml |   


## Apatelet
::: danger  
All environment variables are set by the control plane when spawning Apatelets.
Change these might have unforseen consequences!  
:::  

| Variable name | Description | Default value | 
| --- | --- | --- |
| APATELET_LISTEN_ADDRESS | Listen address for the Apatelet | 0.0.0.0 |
| APATELET_LISTEN_PORT | Listen port for the Apatelet | 8086 |
| APATELET_METRICS_PORT | Listen port for Kubernetes daemon | 10250 |
| APATELET_KUBERNETES_PORT | Listen port for Prometheus metrics | 10255 |
| APATELET_KUBE_CONFIG | Path to the kube config | /tmp/apate/config |
| CI_APATELET_K8S_ADDRESS | CIKubernetesAddress is to add an entry to the /etc/hosts file of an apatelet to ensure it can find the k8s cluster. This can be used to fix bugs when running apate in a DinD |  |
| APATELET_CP_ADDRESS | Address that should be used when connecting to the control plane  | localhost |
| APATELET_CP_PORT | Port that should be used when connecting to the control plane | 8085 |
| APATELET_DISABLE_TAINTS | Determines whether to disable taints on this node or not | false |
| APATELET_ENABLE_DEBUG | Enable extra debug messages | false |