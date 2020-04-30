package cluster

import (
	"io/ioutil"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
)

type KubeConfig []byte


func GetKubeConfig(path string) (KubeConfig, error) {
	return ioutil.ReadFile(filepath.Clean(path))
}

// GetConfigForContext returns a kubernetes client configuration for the context given.
func (k KubeConfig) GetConfig() (*rest.Config, error) {

	config, err := clientcmd.RESTConfigFromKubeConfig(k)
	if err != nil {
		return nil, err
	}

	return config, nil
}