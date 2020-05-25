// Package env provides access to environment variable defaults, definitions & utils
package env

import (
	goEnv "github.com/deanishe/go-env"
	"github.com/pkg/errors"
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
type PullPolicy string

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

// Valid returns true if the pull policy is valid
func (p PullPolicy) Valid() bool {
	return p == AlwaysPull || p == AlwaysLocal || p == PullIfNotLocal
}

// DumpAsKeyValue makes a key=value string array from an environment
func DumpAsKeyValue(env interface{}) ([]string, error) {
	envMap, err := goEnv.Dump(env)
	if err != nil {
		return nil, errors.Wrap(err, "failed to dump controlplane env")
	}

	var envArray []string
	for k, v := range envMap {
		envArray = append(envArray, k+"="+v)
	}

	return envArray, nil
}
