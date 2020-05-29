// Package provider handles the interaction with the virtual kubelet library (and node-cli)
package provider

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/node"
	"strconv"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/condition"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"

	cli "github.com/virtual-kubelet/node-cli"
	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
	baseName   = "apatelet"
)

// Provider implements the node-cli (virtual kubelet) interface for a virtual kubelet provider
type Provider struct {
	Pods      podmanager.PodManager
	Resources *scenario.NodeResources
	Cfg       provider.InitConfig
	NodeInfo  node.NodeInfo
	Store     *store.Store
	Stats     *Stats

	Node       *corev1.Node
	Conditions nodeConditions
}

// CreateProvider creates the node-cli (virtual kubelet) command
func CreateProvider(ctx context.Context, env *env.ApateletEnvironment, res *scenario.NodeResources, k8sPort int, metricsPort int, store *store.Store) (*cli.Command, error) {
	op, err := opts.FromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get options from env")
	}

	name := baseName + "-" + res.UUID.String()
	op.KubeConfigPath = env.KubeConfigLocation
	op.ListenPort = int32(k8sPort)
	op.MetricsAddr = ":" + strconv.Itoa(metricsPort)
	op.Provider = baseName
	op.NodeName = name

	nodeInfo, err := node.NewNodeInfo("apatelet", "agent", name, k8sVersion, res.Selector, metricsPort)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes node info")
	}

	node, err := cli.New(ctx,
		cli.WithProvider(baseName, func(cfg provider.InitConfig) (provider.Provider, error) {
			cfg.DaemonPort = int32(k8sPort)
			return NewProvider(podmanager.New(), NewStats(), res, cfg, nodeInfo, store), nil
		}),
		cli.WithBaseOpts(op),
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create new virtual kubelet provider")
	}

	return node, nil
}

// NewProvider returns the provider but with the vk type instead of our own.
func NewProvider(pods podmanager.PodManager, stats *Stats, resources *scenario.NodeResources, cfg provider.InitConfig, nodeInfo node.NodeInfo, store *store.Store) provider.Provider {
	return &Provider{
		Pods:      pods,
		Resources: resources,
		Cfg:       cfg,
		NodeInfo:  nodeInfo,
		Store:     store,
		Stats:     stats,
		Conditions: nodeConditions{
			ready:              condition.New(true, corev1.NodeReady),
			outOfDisk:          condition.New(false, corev1.NodeOutOfDisk),
			memoryPressure:     condition.New(false, corev1.NodeMemoryPressure),
			diskPressure:       condition.New(false, corev1.NodeDiskPressure),
			networkUnavailable: condition.New(false, corev1.NodeNetworkUnavailable),
			pidPressure:        condition.New(false, corev1.NodePIDPressure),
		},
	}
}
