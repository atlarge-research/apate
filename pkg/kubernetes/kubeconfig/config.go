// Package kubeconfig provides the ability to create, read, and manage the kubeconfig file/bytes.
package kubeconfig

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeConfig contains a bytearray representing the kubeconfig and a path to a file in which this kubeconfig is written.
// It makes sure that this path is only written to once and read from as little as possible.
type KubeConfig struct {
	Path  string
	Bytes []byte
}

// FromBytes creates a kubeConfig struct from byte array and writes it to the given path if the file doesn't exist.
func FromBytes(bytes []byte, path string) (*KubeConfig, error) {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed opening file from provided path")
		}

		if _, err := file.Write(bytes); err != nil {
			return nil, errors.Wrapf(err, "failed to write Kubeconfig to file at %v", file.Name())
		}
	}

	return &KubeConfig{
		Path:  path,
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
