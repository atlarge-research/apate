package cluster

import (
	"io/ioutil"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeConfig is an alias of a bytearray, and represents a raw kube configuration file loaded from file.
type KubeConfig []byte

// GetKubeConfig Loads a KubeConfig from a file path.
func GetKubeConfig(path string) (KubeConfig, error) {
	return ioutil.ReadFile(filepath.Clean(path))
}

// GetConfig returns a kubernetes rest configuration from the KubeConfig.
func (k KubeConfig) GetConfig() (*rest.Config, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(k)
	if err != nil {
		return nil, err
	}

	return config, nil
}
