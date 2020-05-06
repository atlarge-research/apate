// Package container provides methods to create containers for the required Apate components
// and retrieve data from the environment and their defaults
package env

import (
	"os"
	"strconv"
)

const (
	// DefaultPullPolicy returns the default pull policy
	DefaultPullPolicy = PullIfNotLocal

	// AlwaysPull will always pull the image, even if it's already locally available
	AlwaysPull = "pull-always"

	// AlwaysLocal will always use the local image, and will not pull if it's not locally available
	AlwaysLocal = "local-always"

	// PullIfNotLocal will only pull if the image is not locally available, this will not check
	// if the local image is possibly outdated
	PullIfNotLocal = "pull-if-not-local"

	// DockerAddressPrefix specifies the docker address prefix, used for determining the docker address
	DockerAddressPrefix = "172.17."
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

	// ManagedClusterConfig is the path to the config of the cluster manager, if applicable
	ManagedClusterConfig = "CP_K8S_CONFIG"
	// ManagedClusterConfigDefault is the default value for ManagedClusterConfig
	ManagedClusterConfigDefault = "/tmp/apate/manager"

	// ControlPlaneExternalIP can be used to override the IP the control plane will give to apatelets to connect to
	ControlPlaneExternalIP = "CP_EXTERNAL_IP"
	// ControlPlaneExternalIPDefault is the default for ControlPlaneExternalIP
	ControlPlaneExternalIPDefault = "auto"

	// ControlPlaneDockerPolicy specifies the docker pull policy for apatelet images
	ControlPlaneDockerPolicy = "CP_DOCKER_POLICY"
	// ControlPlaneDockerPolicyDefault is the default for ControlPlaneDockerPolicy
	ControlPlaneDockerPolicyDefault = DefaultPullPolicy
)

// Apatelet environment variables
const (
	// ApateletListenAddress is the address the apatelet will listen on for requests
	ApateletListenAddress = "APATELET_LISTEN_ADDRESS"
	// ApateletListenAddressDefault is the default for ApateletListenAddress
	ApateletListenAddressDefault = "0.0.0.0"

	// ApateletListenPort is the port the apatelet will listen on for requests
	ApateletListenPort = "APATELET_LISTEN_PORT"
	// ApateletListenPortDefault is the default for ApateletListenPort
	ApateletListenPortDefault = "8086"

	// ControlPlaneAddress is the address of the control plane which will be used to connect to
	ControlPlaneAddress = "CP_ADDRESS"
	// ControlPlaneAddressDefault is the default for ControlPlaneAddress
	ControlPlaneAddressDefault = "localhost"

	// ControlPlanePort is the port of the control plane
	ControlPlanePort = "CP_PORT"
	// ControlPlanePortDefault is the default for ControlPlanePort
	ControlPlanePortDefault = ControlPlaneListenPortDefault
)

// Docker constants
const (
	// General apate docker constant
	ApateDocker = "apatekubernetes"

	// Apatelet docker constants
	ApateletContainerPrefix = "apatelet-"
	ApateletImageName       = "apatelet:latest"
	ApateletFullImage       = ApateDocker + "/" + ApateletImageName

	// Docker docker constants
	ControlPlaneContainerName = "apate-cp"
	ControlPlaneImageName     = "controlplane:latest"
	ControlPlaneFullImage     = ApateDocker + "/" + ControlPlaneImageName
)

// ApateletEnvironment represents the environment variables of the apatelet
type ApateletEnvironment struct {
	ListenAddress string
	ListenPort    int

	ControlPlaneAddress string
	ControlPlanePort    int
}

// DefaultApateletEnvironment returns the default apate environment
func DefaultApateletEnvironment() ApateletEnvironment {
	defaultPort, _ := strconv.Atoi(ApateletListenPortDefault)
	return ApateletEnvironment{
		ListenAddress: ApateletListenAddressDefault,
		ListenPort:    defaultPort,
	}
}

func ApateletEnvironmentFromEnv() (ApateletEnvironment, error) {
	controlPlaneAddress := RetrieveFromEnvironment(ControlPlaneAddress, ControlPlaneAddressDefault)

	controlPlanePort, err := strconv.Atoi(RetrieveFromEnvironment(ControlPlanePort, ControlPlanePortDefault))
	if err != nil {
		return ApateletEnvironment{}, err
	}

	// Retrieve own port
	listenPort, err := strconv.Atoi(RetrieveFromEnvironment(ApateletListenPort, ApateletListenPortDefault))
	if err != nil {
		return ApateletEnvironment{}, err
	}

	// Retrieving connection information
	listenAddress := RetrieveFromEnvironment(ApateletListenAddress, ApateletListenAddressDefault)

	return ApateletEnvironment{
		ListenAddress:       listenAddress,
		ListenPort:          listenPort,
		ControlPlaneAddress: controlPlaneAddress,
		ControlPlanePort:    controlPlanePort,
	}, nil
}

func (env *ApateletEnvironment) AddConnectionInfo(address string, port int) {
	env.ControlPlaneAddress = address
	env.ControlPlanePort = port
}

func (env *ApateletEnvironment) Copy() ApateletEnvironment {
	return *env
}

// RetrieveFromEnvironment allows for a value to be retrieved from the environment
func RetrieveFromEnvironment(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}
