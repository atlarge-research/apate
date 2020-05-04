// Package container provides methods to create containers for the required Apate components
// and retrieve data from the environment and their defaults
package container

import (
	"os"
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
	apateDocker = "apatekubernetes"

	// Apatelet docker constants
	apateletContainerPrefix = "apatelet-"
	apateletImageName       = "apatelet:latest"
	apateletFullImage       = apateDocker + "/" + apateletImageName

	// Docker docker constants
	controlPlaneContainerName = "apate-cp"
	controlPlaneImageName     = "controlplane:latest"
	controlPlaneFullImage     = apateDocker + "/" + controlPlaneImageName
)

// GetPullPolicyControlPlane returns the pull policy active on the control plane
func GetPullPolicyControlPlane() string {
	var policy string
	if val, ok := os.LookupEnv(ControlPlaneDockerPolicy); ok {
		policy = val
	} else {
		policy = DefaultPullPolicy
	}

	return policy
}

// RetrieveFromEnvironment allows for a value to be retrieved from the environment
func RetrieveFromEnvironment(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}
