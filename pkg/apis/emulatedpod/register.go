package emulatedpod

import (
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster/kubeconfig"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubectl"
	"io/ioutil"
)

const (
	GroupName = "apate.opendc.org"
)

func AddCRDToKubernetes(config *kubeconfig.KubeConfig) error {

	file, err := ioutil.ReadFile("config/crd/apate.opendc.org_emulatedpods.yaml")
	if err != nil {
		return err
	}

	if err := kubectl.Apply(file, config); err != nil {
		return err
	}

	return nil
}
