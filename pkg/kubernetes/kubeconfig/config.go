// Package kubeconfig provides the ability to create, read, and manage the kubeconfig file/bytes.
package kubeconfig

import (
	"io/ioutil"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeConfig is an alias of a bytearray, and represents a raw kube configuration file loaded from file.
type KubeConfig struct {
	Path  string
	Bytes []byte
}

// FromBytes creates a kubeConfig struct from byte array.
func FromBytes(bytes []byte) (*KubeConfig, error) {
	file, err := ioutil.TempFile("", "config")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create tmp file for kubeconfig")
	}

	if _, err := file.Write(bytes); err != nil {
		return nil, errors.Wrapf(err, "failed to write Kubeconfig to file at %v", file.Name())
	}

	return &KubeConfig{
		Path:  file.Name(),
		Bytes: bytes,
	}, nil
}

// FromPath Loads a KubeConfig from a file path.
func FromPath(path string) (*KubeConfig, error) {
	bytes, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read kubeconfig from file at %v", path)
	}

	return &KubeConfig{path, bytes}, nil
}

// GetConfig returns a kubernetes rest configuration from the KubeConfig.
func (k KubeConfig) GetConfig() (*rest.Config, error) {
	config, err := clientcmd.RESTConfigFromKubeConfig(k.Bytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create rest config")
	}

	return config, nil
}
