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

// ApateletEnvironment represents the environment variables of the apatelet
type ApateletEnvironment struct {
	ListenAddress string
	ListenPort    int

	ControlPlaneAddress string
	ControlPlanePort    int
}

// DefaultApateletEnvironment returns the default apate environment
func DefaultApateletEnvironment() (ApateletEnvironment, error) {
	defaultPort, err := strconv.Atoi(ApateletListenPortDefault)
	if err != nil {
		return ApateletEnvironment{}, errors.Wrapf(err, "failed to convert Apatelet listening port (%v) to string", ApateletListenPortDefault)
	}

	return ApateletEnvironment{
		ListenAddress: ApateletListenAddressDefault,
		ListenPort:    defaultPort,
	}, nil
}

// ApateletEnvironmentFromEnv build an ApateletEnvironment based on the actual environment
func ApateletEnvironmentFromEnv() (ApateletEnvironment, error) {
	controlPlaneAddress := RetrieveFromEnvironment(ControlPlaneAddress, ControlPlaneAddressDefault)

	cpp := RetrieveFromEnvironment(ControlPlanePort, ControlPlanePortDefault)
	controlPlanePort, err := strconv.Atoi(cpp)
	if err != nil {
		return ApateletEnvironment{}, errors.Wrapf(err, "failed to convert contolplane port (%v) to string", cpp)
	}

	// Retrieve own port
	lp := RetrieveFromEnvironment(ApateletListenPort, ApateletListenPortDefault)
	listenPort, err := strconv.Atoi(lp)
	if err != nil {
		return ApateletEnvironment{}, errors.Wrapf(err, "failed to convert default Apatelet listening port (%v) to string", lp)
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

// AddConnectionInfo adds the given connection information to the environment
func (env *ApateletEnvironment) AddConnectionInfo(address string, port int) {
	env.ControlPlaneAddress = address
	env.ControlPlanePort = port
}
