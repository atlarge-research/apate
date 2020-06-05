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

	// ApateletControlPlaneAddressDefault is the default for ControlPlaneAddress
	ApateletControlPlaneAddressDefault = "localhost"

	// ApateletControlPlanePortDefault is the default for ControlPlanePort
	ApateletControlPlanePortDefault = CPListenPortDefault

	// ApateletDisableTaintsDefault is the default for DisableTaints
	ApateletDisableTaintsDefault = false
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
	KubernetesAddress  string `env:"APATELET_K8S_ADDRESS"`

	// ControlPlaneAddress is the address of the control plane which will be used to connect to
	ControlPlaneAddress string `env:"APATELET_CP_ADDRESS"`
	// ControlPlanePort is the port of the control plane
	ControlPlanePort int `env:"APATELET_CP_PORT"`

	// DisableTaints determines whether to disable taints on this node or not
	DisableTaints bool `env:"APATELET_DISABLE_TAINTS"`
}

// defaultApateletEnvironment returns the default apate environment
func defaultApateletEnvironment() ApateletEnvironment {
	return ApateletEnvironment{
		ListenAddress: ApateletListenAddressDefault,
		ListenPort:    ApateletListenPortDefault,

		KubernetesPort: ApateletKubernetesPortDefault,
		MetricsPort:    ApateletMetricsPortDefault,

		KubeConfigLocation: ApateletKubeConfigLocationDefault,
		KubernetesAddress:  "",

		ControlPlaneAddress: ApateletControlPlaneAddressDefault,
		ControlPlanePort:    ApateletControlPlanePortDefault,

		DisableTaints: ApateletDisableTaintsDefault,
	}
}

// ApateletEnv builds an ApateletEnvironment based on the actual environment
func ApateletEnv() (ApateletEnvironment, error) {
	c := defaultApateletEnvironment()
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
