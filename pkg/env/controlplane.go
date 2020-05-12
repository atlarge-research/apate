package env

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
type RunType = string

const (
	// Routine uses go routines for the creation of apatelets
	Routine RunType = "ROUTINES"

	// Docker uses docker containers for the creation of apatelets
	Docker RunType = "DOCKER"
)

// ControlPlaneEnvironment represents the environment variables of the control plane
type ControlPlaneEnvironment struct {
	Address, Port, ManagerConfig, ExternalIP, DockerPolicy, ApateletRunType, PrometheusStackEnabled string
}

// DefaultControlPlaneEnvironment returns the default control plane environment
func DefaultControlPlaneEnvironment() ControlPlaneEnvironment {
	return ControlPlaneEnvironment{
		Address:                ControlPlaneListenAddressDefault,
		Port:                   ControlPlaneListenPortDefault,
		ManagerConfig:          ManagedClusterConfigDefault,
		ExternalIP:             ControlPlaneExternalIPDefault,
		DockerPolicy:           ControlPlaneDockerPolicyDefault,
		ApateletRunType:        ControlPlaneApateletRunTypeDefault,
		PrometheusStackEnabled: PrometheusStackEnabledDefault,
	}
}
