package env

import (
	"strconv"

	"github.com/pkg/errors"
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
	ApateletListenPortDefault = 8086

	// ApateletKubeConfigLocation is the path to the kube config
	ApateletKubeConfigLocation = "APATELET_KUBE_CONFIG"
	// ApateletKubeConfigLocationDefault is the default value for ApateletKubeConfigLocation
	ApateletKubeConfigLocationDefault = "/tmp/apate/config"

	// ControlPlaneAddress is the address of the control plane which will be used to connect to
	ControlPlaneAddress = "CP_ADDRESS"
	// ControlPlaneAddressDefault is the default for ControlPlaneAddress
	ControlPlaneAddressDefault = "localhost"

	// ControlPlanePort is the port of the control plane
	ControlPlanePort = "CP_PORT"
	// ControlPlanePortDefault is the default for ControlPlanePort
	ControlPlanePortDefault = ControlPlaneListenPortDefault
)

// ApateletEnvironment represents the environment variables of the apatelet
type ApateletEnvironment struct {
	ListenAddress string
	ListenPort    int

	KubeConfigLocation string

	ControlPlaneAddress string
	ControlPlanePort    int
}

// DefaultApateletEnvironment returns the default apate environment
func DefaultApateletEnvironment() ApateletEnvironment {
	return ApateletEnvironment{
		ListenAddress: ApateletListenAddressDefault,
		ListenPort:    ApateletListenPortDefault,

		KubeConfigLocation: ApateletKubeConfigLocationDefault,

		ControlPlaneAddress: ControlPlaneAddressDefault,
		ControlPlanePort:    ControlPlanePortDefault,
	}
}

// ApateletEnv builds an ApateletEnvironment based on the actual environment
func ApateletEnv() (ApateletEnvironment, error) {
	cpp := RetrieveFromEnvironment(ControlPlanePort, strconv.Itoa(ControlPlanePortDefault))
	controlPlanePort, err := strconv.Atoi(cpp)
	if err != nil {
		return ApateletEnvironment{}, errors.Wrapf(err, "invalid port %v while parsing controlPlanePort", cpp)
	}

	// Retrieve own port
	lp := RetrieveFromEnvironment(ApateletListenPort, strconv.Itoa(ApateletListenPortDefault))
	listenPort, err := strconv.Atoi(lp)
	if err != nil {
		return ApateletEnvironment{}, errors.Wrapf(err, "invalid port %v while parsing listenPort", lp)
	}

	return ApateletEnvironment{
		ListenAddress:       RetrieveFromEnvironment(ApateletListenAddress, ApateletListenAddressDefault),
		ListenPort:          listenPort,
		KubeConfigLocation:  RetrieveFromEnvironment(ApateletKubeConfigLocation, ApateletKubeConfigLocationDefault),
		ControlPlaneAddress: RetrieveFromEnvironment(ControlPlaneAddress, ControlPlaneAddressDefault),
		ControlPlanePort:    controlPlanePort,
	}, nil
}

// AddConnectionInfo adds the given connection information to the environment
func (env *ApateletEnvironment) AddConnectionInfo(address string, port int) {
	env.ControlPlaneAddress = address
	env.ControlPlanePort = port
}
