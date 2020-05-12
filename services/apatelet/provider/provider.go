// Package provider handles the interaction with the virtual kubelet library (and node-cli)
package provider

import (
	"context"
	"os"
	"strconv"

	emulatedpodv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/emulatedpod/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"

	cli "github.com/virtual-kubelet/node-cli"
	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
	baseName   = "apatelet"
)

// Provider implements the node-cli (virtual kubelet) interface for a virtual kubelet provider
type Provider struct {
	pods        podmanager.PodManager
	resources   *normalization.NodeResources
	cfg         provider.InitConfig
	nodeInfo    cluster.NodeInfo
	store       *store.Store
	crdInformer *emulatedpodv1.Informer
}

// CreateProvider creates the node-cli (virtual kubelet) command
func CreateProvider(ctx context.Context, res *normalization.NodeResources, k8sPort int, metricsPort int, store *store.Store, crdInformer *emulatedpodv1.Informer) (*cli.Command, error) {
	op, err := opts.FromEnv()
	if err != nil {
		return nil, err
	}

	name := baseName + "-" + res.UUID.String()
	op.KubeConfigPath = os.TempDir() + "/apate/config"
	op.ListenPort = int32(k8sPort)
	op.MetricsAddr = ":" + strconv.Itoa(metricsPort)
	op.Provider = baseName
	op.NodeName = name

	nodeInfo := cluster.NewNode("virtual-kubelet", "agent", name, k8sVersion)

	node, err := cli.New(ctx,
		cli.WithProvider(baseName, func(cfg provider.InitConfig) (provider.Provider, error) {
			cfg.DaemonPort = int32(k8sPort)
			return NewProvider(res, cfg, nodeInfo, store, crdInformer), nil
		}),
		cli.WithBaseOpts(op),
	)

	return node, err
}

// NewProvider returns the provider but with the vk type instead of our own.
func NewProvider(resources *normalization.NodeResources, cfg provider.InitConfig, nodeInfo cluster.NodeInfo, store *store.Store, crdInformer *emulatedpodv1.Informer) provider.Provider {
	return &Provider{
		pods:        podmanager.New(),
		resources:   resources,
		cfg:         cfg,
		nodeInfo:    nodeInfo,
		store:       store,
		crdInformer: crdInformer,
	}
}
