package env

import (
	"log"

	"github.com/deanishe/go-env"
	"github.com/pkg/errors"
)

// Control plane environment variables
const (
	// ControlPlaneListenAddressDefault is the default value for ControlPlaneListenAddress
	ControlPlaneListenAddressDefault = "0.0.0.0"

	// ControlPlaneListenPortDefault is the default value for ControlPlaneListenPort
	ControlPlaneListenPortDefault = 8085

	// ManagedClusterConfigDefault is the default value for ManagedClusterConfigLocation
	ManagedClusterConfigLocationDefault = "/tmp/apate/manager"

	// ControlPlaneKubeConfigLocationDefault is the default value for ControlPlaneKubeConfigLocation
	ControlPlaneKubeConfigLocationDefault = "/tmp/apate/config"

	// ControlPlaneExternalIPDefault is the default for ControlPlaneExternalIP
	ControlPlaneExternalIPDefault = "auto"

	// ControlPlaneDockerPolicyDefault is the default for ControlPlaneDockerPolicy
	ControlPlaneDockerPolicyDefault = DefaultPullPolicy

	// ControlPlaneApateletRunTypeDefault is the default for ControlPlaneApateletRunType
	ControlPlaneApateletRunTypeDefault = Routine

	// PrometheusStackEnabledDefault is the default for PrometheusStackEnabled
	PrometheusStackEnabledDefault = true

	// Pod/Node CRD Defaults
	NodeCRDLocationDefault = "config/crd/apate.opendc.org_nodeconfigurations.yaml"
	PodCRDLocationDefault  = "config/crd/apate.opendc.org_podconfigurations.yaml"

	// KinD default cluster name
	KinDClusterNameDefault = "Apate"
)

// RunType is the run strategy used by the control plane to run apalets
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
	ManagerConfigLocation string `env:"CP_MANAGER_LOCATION"`
	// KubeConfigLocation is the path to the kube config
	KubeConfigLocation string `env:"CP_KUBE_CONFIG"`

	// DockerPolicy specifies the docker pull policy for apatelet images
	DockerPolicy PullPolicy `env:"CP_DOCKER_POLICY"`
	// ApateletRunType specifies how the control plane runs new apatelets
	ApateletRunType RunType `env:"CP_APATELET_RUN_TYPE"`

	// PrometheusStackEnabled specifies
	PrometheusStackEnabled bool `env:"CP_PROMETHEUS"`

	// CRD Locations
	PodCRDLocation  string `env:"CP_POD_CRD_LOCATION"`
	NodeCRDLocation string `env:"CP_NODE_CRD_LOCATION"`

	// (KinD) Cluster Name
	KinDClusterName string `env:"KIND_CLUSTER_NAME"`
}

var controlPlaneEnvironment *ControlPlaneEnvironment

// DefaultControlPlaneEnvironment returns the default control plane environment
func DefaultControlPlaneEnvironment() ControlPlaneEnvironment {
	return ControlPlaneEnvironment{
		ListenAddress:          ControlPlaneListenAddressDefault,
		ListenPort:             ControlPlaneListenPortDefault,
		ManagerConfigLocation:  ManagedClusterConfigLocationDefault,
		KubeConfigLocation:     ControlPlaneKubeConfigLocationDefault,
		ExternalIP:             ControlPlaneExternalIPDefault,
		DockerPolicy:           ControlPlaneDockerPolicyDefault,
		ApateletRunType:        ControlPlaneApateletRunTypeDefault,
		PrometheusStackEnabled: PrometheusStackEnabledDefault,
		NodeCRDLocation:        NodeCRDLocationDefault,
		PodCRDLocation:         PodCRDLocationDefault,
		KinDClusterName:        KinDClusterNameDefault,
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
