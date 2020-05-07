// Package provider handles the interaction with the virtual kubelet library (and node-cli)
package provider

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/cluster"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/normalization"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	cli "github.com/virtual-kubelet/node-cli"
	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"
	"os"
	"strconv"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
	baseName   = "apatelet"
)

// Provider implements the node-cli (virtual kubelet) interface for a virtual kubelet provider
type Provider struct {
	pods      PodManager
	resources *normalization.NodeResources
	cfg       provider.InitConfig
	nodeInfo  cluster.NodeInfo
	store     *store.Store
}

// CreateProvider creates the node-cli (virtual kubelet) command
func CreateProvider(ctx context.Context, res *normalization.NodeResources, k8sPort int, metricsPort int, store *store.Store) (*cli.Command, error) {
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
			return NewProvider(res, cfg, nodeInfo, store), nil
		}),
		cli.WithBaseOpts(op),
	)

	return node, err
}

// NewProvider returns the provider but with the vk type instead of our own.
func NewProvider(resources *normalization.NodeResources, cfg provider.InitConfig, nodeInfo cluster.NodeInfo, store *store.Store) provider.Provider {
	return &Provider{
		pods:      NewPodManager(),
		resources: resources,
		cfg:       cfg,
		nodeInfo:  nodeInfo,
		store:     store,
	}
}
