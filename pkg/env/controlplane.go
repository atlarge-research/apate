package env

import (
	"strconv"

	"github.com/pkg/errors"
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
	ControlPlaneListenPortDefault = "8085"

	// ManagedClusterConfigLocation is the path to the config of the cluster manager, if applicable
	ManagedClusterConfigLocation = "CP_MANAGER_LOCATION"
	// ManagedClusterConfigDefault is the default value for ManagedClusterConfigLocation
	ManagedClusterConfigLocationDefault = "/tmp/apate/manager"

	// KubeConfigLocation is the path to the kube config
	KubeConfigLocation = "CP_KUBE_CONFIG"
	// KubeConfigLocationDefault is the default value for KubeConfigLocation
	KubeConfigLocationDefault = "/tmp/apate/config"

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
	PrometheusStackEnabledDefault = "true"
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
func DefaultControlPlaneEnvironment() (ControlPlaneEnvironment, error) {
	defaultPort, err := strconv.Atoi(ControlPlaneListenPortDefault)
	if err != nil {
		return ControlPlaneEnvironment{}, errors.Wrapf(err, "failed to convert Apatelet listening port (%v) to string", ApateletListenPortDefault)
	}

	prometheusEnabled, err := strconv.ParseBool(PrometheusStackEnabledDefault)
	if err != nil {
		return ControlPlaneEnvironment{}, errors.Wrapf(err, "failed to convert Apatelet listening port (%v) to string", ApateletListenPortDefault)
	}

	return ControlPlaneEnvironment{
		ListenAddress:          ControlPlaneListenAddressDefault,
		ListenPort:             defaultPort,
		ManagerConfigLocation:  ManagedClusterConfigLocationDefault,
		KubeConfigLocation:     KubeConfigLocationDefault,
		ExternalIP:             ControlPlaneExternalIPDefault,
		DockerPolicy:           ControlPlaneDockerPolicyDefault,
		ApateletRunType:        ControlPlaneApateletRunTypeDefault,
		PrometheusStackEnabled: prometheusEnabled,
	}, nil
}

// SetEnv overrides the current environment for the control plane
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
		panic("invalid pull policy: " + pullPolicy)
	}

	port := RetrieveFromEnvironment(ControlPlaneListenPort, ControlPlaneListenPortDefault)
	listenPort, err := strconv.Atoi(port)
	if err != nil {
		panic("invalid port: " + err.Error())
	}

	prometheus := RetrieveFromEnvironment(PrometheusStackEnabled, PrometheusStackEnabledDefault)
	prometheusEnabled, err := strconv.ParseBool(prometheus)
	if err != nil {
		panic("invalid boolean for prometheus enabled: " + port)
	}

	controlPlaneEnvironment = &ControlPlaneEnvironment{
		ListenAddress:          RetrieveFromEnvironment(ControlPlaneListenAddress, ControlPlaneListenAddressDefault),
		ListenPort:             listenPort,
		ManagerConfigLocation:  RetrieveFromEnvironment(ManagedClusterConfigLocation, ManagedClusterConfigLocationDefault),
		ExternalIP:             RetrieveFromEnvironment(ControlPlaneExternalIP, ControlPlaneExternalIPDefault),
		KubeConfigLocation:     RetrieveFromEnvironment(KubeConfigLocation, KubeConfigLocationDefault),
		DockerPolicy:           pullPolicy,
		ApateletRunType:        RunType(RetrieveFromEnvironment(ControlPlaneApateletRunType, string(ControlPlaneApateletRunTypeDefault))),
		PrometheusStackEnabled: prometheusEnabled,
	}

	return *controlPlaneEnvironment
}
