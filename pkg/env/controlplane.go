package env

import (
	"log"

	"github.com/deanishe/go-env"
	"github.com/pkg/errors"
)

// Control plane environment variables
const (
	// CPListenAddressDefault is the default value for ControlPlaneListenAddress
	CPListenAddressDefault = "0.0.0.0"

	// CPListenPortDefault is the default value for ControlPlaneListenPort
	CPListenPortDefault = 8085

	// CPManagedClusterConfigLocationDefault is the default value for ManagedClusterConfigLocation
	CPManagedClusterConfigLocationDefault = "config/kind.yml"

	// CPKubeConfigLocationDefault is the default value for ControlPlaneKubeConfigLocation
	CPKubeConfigLocationDefault = "/tmp/apate/config"

	// CPExternalIPDefault is the default for ControlPlaneExternalIP
	CPExternalIPDefault = "auto"

	// CPDockerPolicyDefault is the default for ControlPlaneDockerPolicy
	CPDockerPolicyDefault = DefaultPullPolicy

	// CPApateletRunTypeDefault is the default for ControlPlaneApateletRunType
	CPApateletRunTypeDefault = Routine

	// CPPrometheusStackEnabledDefault is the default for PrometheusStackEnabled
	CPPrometheusStackEnabledDefault = true

	// CPNodeCRDLocationDefault CRD default location
	CPNodeCRDLocationDefault = "config/crd/apate.opendc.org_nodeconfigurations.yaml"
	// CPPodCRDLocationDefault CRD default location
	CPPodCRDLocationDefault = "config/crd/apate.opendc.org_podconfigurations.yaml"

	// CPKinDClusterNameDefault default cluster name
	CPKinDClusterNameDefault = "apate"

	// CPUseKinDInternalConfig default for UseDockerHostname
	CPUseDockerHostnameDefault = false

	// CPDebugEnabledDefault default for DebugEnabled
	CPDebugEnabledDefault = false
)

// RunType is the runner strategy used by the control plane to run apalets
type RunType string

const (
	// Routine uses go routines for the creation of apatelets
	Routine RunType = "ROUTINES"

	// Docker uses docker containers for the creation of apatelets
	Docker RunType = "DOCKER"
)

// ControlPlaneEnvironment represents the environment variables of the control plane
type ControlPlaneEnvironment struct {
	// ListenAddress is the address the control plane will listen on
	ListenAddress string `env:"CP_LISTEN_ADDRESS"`
	// ListenPort is the port the control plane will listen on
	ListenPort int `env:"CP_LISTEN_PORT"`

	// ExternalIP can be used to override the IP the control plane will give to apatelets to connect to
	ExternalIP string `env:"CP_EXTERNAL_IP"`

	// ManagerConfigLocation is the path to the config of the cluster manager, if applicable
	ManagerConfigLocation string `env:"CP_MANAGER_CONFIG_LOCATION"`
	// KubeConfigLocation is the path to the kube config
	KubeConfigLocation string `env:"CP_KUBE_CONFIG_LOCATION"`

	// DockerPolicy specifies the docker pull policy for apatelet images
	DockerPolicy PullPolicy `env:"CP_DOCKER_POLICY"`
	// ApateletRunType specifies how the control plane runs new apatelets
	ApateletRunType RunType `env:"CP_APATELET_RUN_TYPE"`

	// PrometheusStackEnabled specifies if the control plane should create a prometheus stack on startup
	PrometheusStackEnabled bool `env:"CP_PROMETHEUS"`

	// CRD Locations
	PodCRDLocation  string `env:"CP_POD_CRD_LOCATION"`
	NodeCRDLocation string `env:"CP_NODE_CRD_LOCATION"`

	// (KinD) Cluster Name
	KinDClusterName string `env:"CP_KIND_CLUSTER_NAME"`

	// UseDockerHostname specifies if we should rewrite the KinD address to 'docker'
	UseDockerHostname bool `env:"CP_DOCKER_HOSTNAME"`

	DebugEnabled bool `env:"CP_ENABLE_DEBUG"`
}

var controlPlaneEnvironment *ControlPlaneEnvironment

// DefaultControlPlaneEnvironment returns the default control plane environment
func DefaultControlPlaneEnvironment() ControlPlaneEnvironment {
	return ControlPlaneEnvironment{
		ListenAddress:          CPListenAddressDefault,
		ListenPort:             CPListenPortDefault,
		ManagerConfigLocation:  CPManagedClusterConfigLocationDefault,
		KubeConfigLocation:     CPKubeConfigLocationDefault,
		ExternalIP:             CPExternalIPDefault,
		DockerPolicy:           CPDockerPolicyDefault,
		ApateletRunType:        CPApateletRunTypeDefault,
		PrometheusStackEnabled: CPPrometheusStackEnabledDefault,
		NodeCRDLocation:        CPNodeCRDLocationDefault,
		PodCRDLocation:         CPPodCRDLocationDefault,
		KinDClusterName:        CPKinDClusterNameDefault,
		UseDockerHostname:      CPUseDockerHostnameDefault,
		DebugEnabled:           CPDebugEnabledDefault,
	}
}

// SetEnv overrides the current environment for the control plane
// We preferred this over providing a pointer in the getter to avoid accidental overrides
func SetEnv(environment ControlPlaneEnvironment) {
	controlPlaneEnvironment = &environment
}

// ControlPlaneEnv builds and ControlPlaneEnvironment based on the actual environment
func ControlPlaneEnv() ControlPlaneEnvironment {
	if controlPlaneEnvironment != nil {
		return *controlPlaneEnvironment
	}

	c := DefaultControlPlaneEnvironment()
	if err := env.Bind(&c); err != nil {
		log.Panicf("%+v", errors.Wrap(err, "unable to bind control plane environment"))
	}

	controlPlaneEnvironment = &c
	return *controlPlaneEnvironment
}
