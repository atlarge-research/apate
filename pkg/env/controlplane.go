package env

import (
	"log"
	"strconv"
)

// Control plane environment variables
const (
	// ControlPlaneListenAddress is the address the control plane will listen on
	ControlPlaneListenAddress = "CP_LISTEN_ADDRESS"
	// ControlPlaneListenAddressDefault is the default value for ControlPlaneListenAddress
	ControlPlaneListenAddressDefault = "0.0.0.0"

	// ControlPlaneListenPort is the port the control plane will listen on
	ControlPlaneListenPort = "CP_LISTEN_PORT"
	// ControlPlaneListenPortDefault is the default value for ControlPlaneListenPort
	ControlPlaneListenPortDefault = 8085

	// ManagedClusterConfigLocation is the path to the config of the cluster manager, if applicable
	ManagedClusterConfigLocation = "CP_MANAGER_LOCATION"
	// ManagedClusterConfigDefault is the default value for ManagedClusterConfigLocation
	ManagedClusterConfigLocationDefault = "/tmp/apate/manager"

	// ControlPlaneKubeConfigLocation is the path to the kube config
	ControlPlaneKubeConfigLocation = "CP_KUBE_CONFIG"
	// ControlPlaneKubeConfigLocationDefault is the default value for ControlPlaneKubeConfigLocation
	ControlPlaneKubeConfigLocationDefault = "/tmp/apate/config"

	// ControlPlaneExternalIP can be used to override the IP the control plane will give to apatelets to connect to
	ControlPlaneExternalIP = "CP_EXTERNAL_IP"
	// ControlPlaneExternalIPDefault is the default for ControlPlaneExternalIP
	ControlPlaneExternalIPDefault = "auto"

	// ControlPlaneDockerPolicy specifies the docker pull policy for apatelet images
	ControlPlaneDockerPolicy = "CP_DOCKER_POLICY"
	// ControlPlaneDockerPolicyDefault is the default for ControlPlaneDockerPolicy
	ControlPlaneDockerPolicyDefault = DefaultPullPolicy

	// ControlPlaneApateletRunType specifies how the control plane runs new apatelets
	ControlPlaneApateletRunType = "CP_APATELET_RUN_TYPE"
	// ControlPlaneApateletRunTypeDefault is the default for ControlPlaneApateletRunType
	ControlPlaneApateletRunTypeDefault = Routine

	// PrometheusStackEnabled specifies
	PrometheusStackEnabled = "CP_PROMETHEUS"
	// PrometheusStackEnabledDefault is the default for PrometheusStackEnabled
	PrometheusStackEnabledDefault = true
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
	ListenAddress, ExternalIP string
	ListenPort                int

	ManagerConfigLocation, KubeConfigLocation string

	DockerPolicy    PullPolicy
	ApateletRunType RunType

	PrometheusStackEnabled bool
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

	pullPolicy := PullPolicy(RetrieveFromEnvironment(ControlPlaneDockerPolicy, string(ControlPlaneDockerPolicyDefault)))
	if !pullPolicy.Valid() {
		log.Panicf("invalid pull policy %v when creating control plane env", pullPolicy)
	}

	port := RetrieveFromEnvironment(ControlPlaneListenPort, strconv.Itoa(ControlPlaneListenPortDefault))
	listenPort, err := strconv.Atoi(port)
	if err != nil {
		log.Panicf("invalid port %v when parsing listen port %+v to create control plane env", port, err)
	}

	prometheus := RetrieveFromEnvironment(PrometheusStackEnabled, strconv.FormatBool(PrometheusStackEnabledDefault))
	prometheusEnabled, err := strconv.ParseBool(prometheus)
	if err != nil {
		log.Panicf("invalid boolean %v for prometheus enabled %v to create control plane env", prometheus, port)
	}

	controlPlaneEnvironment = &ControlPlaneEnvironment{
		ListenAddress:          RetrieveFromEnvironment(ControlPlaneListenAddress, ControlPlaneListenAddressDefault),
		ListenPort:             listenPort,
		ManagerConfigLocation:  RetrieveFromEnvironment(ManagedClusterConfigLocation, ManagedClusterConfigLocationDefault),
		ExternalIP:             RetrieveFromEnvironment(ControlPlaneExternalIP, ControlPlaneExternalIPDefault),
		KubeConfigLocation:     RetrieveFromEnvironment(ControlPlaneKubeConfigLocation, ControlPlaneKubeConfigLocationDefault),
		DockerPolicy:           pullPolicy,
		ApateletRunType:        RunType(RetrieveFromEnvironment(ControlPlaneApateletRunType, string(ControlPlaneApateletRunTypeDefault))),
		PrometheusStackEnabled: prometheusEnabled,
	}

	return *controlPlaneEnvironment
}
