// Package env provides access to environment variable defaults, definitions & utils
package env

import (
	"os"
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

	// DockerAddressPrefix specifies the docker address prefix, used for determining the docker address
	DockerAddressPrefix = "172.17."
)

// PullPolicy defines the pull policy used for
type PullPolicy = string

const (
	// DefaultPullPolicy returns the default pull policy
	DefaultPullPolicy = PullIfNotLocal

	// AlwaysPull will always pull the image, even if it's already locally available
	AlwaysPull PullPolicy = "pull-always"

	// AlwaysLocal will always use the local image, and will not pull if it's not locally available
	AlwaysLocal PullPolicy = "local-always"

	// PullIfNotLocal will only pull if the image is not locally available, this will not check
	// if the local image is possibly outdated
	PullIfNotLocal PullPolicy = "pull-if-not-local"
)

// RetrieveFromEnvironment allows for a value to be retrieved from the environment
func RetrieveFromEnvironment(key, def string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}

	return def
}
