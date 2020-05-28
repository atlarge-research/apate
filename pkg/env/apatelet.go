package env

import (
	"github.com/deanishe/go-env"
	"github.com/pkg/errors"
)

// Apatelet environment variables
const (
	// ApateletListenAddressDefault is the default for ListenAddress
	ApateletListenAddressDefault = "0.0.0.0"

	// ApateletListenPortDefault is the default for ListenPort
	ApateletListenPortDefault = 8086

	// ApateletKubernetesPortDefault is the default for MetricsPort
	ApateletKubernetesPortDefault = 10250

	// ApateletMetricsPortDefault is the default for KubernetesPort
	ApateletMetricsPortDefault = 10255

	// ApateletKubeConfigLocationDefault is the default value for KubeConfigLocation
	ApateletKubeConfigLocationDefault = "/tmp/apate/config"

	// ControlPlaneAddressDefault is the default for ControlPlaneAddress
	ControlPlaneAddressDefault = "localhost"

	// ControlPlanePortDefault is the default for ControlPlanePort
	ControlPlanePortDefault = CPListenPortDefault
)

// ApateletEnvironment represents the environment variables of the apatelet
type ApateletEnvironment struct {
	// ListenAddress is the address the apatelet will listen on for requests
	ListenAddress string `env:"APATELET_LISTEN_ADDRESS"`
	// ListenPort is the port the apatelet will listen on for gRPC requests
	ListenPort int `env:"APATELET_LISTEN_PORT"`

	// MetricsPort is the port the apatelet will serve metrics on
	MetricsPort int `env:"APATELET_METRICS_PORT"`
	// KubernetesPort is the port the apatelet will create Kubernetes endpoints on
	KubernetesPort int `env:"APATELET_KUBERNETES_PORT"`

	// KubeConfigLocation is the path to the kube config
	KubeConfigLocation string `env:"APATELET_KUBE_CONFIG"`

	// ControlPlaneAddress is the address of the control plane which will be used to connect to
	ControlPlaneAddress string `env:"CP_ADDRESS"`
	// ControlPlanePort is the port of the control plane
	ControlPlanePort int `env:"CP_PORT"`
}

// DefaultApateletEnvironment returns the default apate environment
func DefaultApateletEnvironment() ApateletEnvironment {
	return ApateletEnvironment{
		ListenAddress: ApateletListenAddressDefault,
		ListenPort:    ApateletListenPortDefault,

		KubernetesPort: ApateletKubernetesPortDefault,
		MetricsPort:    ApateletMetricsPortDefault,

		KubeConfigLocation: ApateletKubeConfigLocationDefault,

		ControlPlaneAddress: ControlPlaneAddressDefault,
		ControlPlanePort:    ControlPlanePortDefault,
	}
}

// ApateletEnv builds an ApateletEnvironment based on the actual environment
func ApateletEnv() (ApateletEnvironment, error) {
	c := DefaultApateletEnvironment()
	if err := env.Bind(&c); err != nil {
		return ApateletEnvironment{}, errors.Wrap(err, "invalid environment variables")
	}

	return c, nil
}

// AddConnectionInfo adds the given connection information to the environment
func (env *ApateletEnvironment) AddConnectionInfo(address string, port int) {
	env.ControlPlaneAddress = address
	env.ControlPlanePort = port
}
