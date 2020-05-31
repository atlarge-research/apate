// Package provider handles the interaction with the virtual kubelet library (and node-cli)
package provider

import (
	"context"
	root "github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/commands"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/opts"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/provider"
	"strconv"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/condition"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
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
	NodeInfo  kubernetes.NodeInfo
	Store     *store.Store
	Stats     *Stats

	Node       *corev1.Node
	Conditions nodeConditions
}

// TODO: Move
type VK struct {
	st   *provider.Store
	opts *opts.Opts
}

func (vk *VK) Run(originalCtx context.Context, ctx context.Context) (int, int, error) {
	metricsPort, k8sPort, err := root.RunRootCommand(originalCtx, ctx, vk.st, vk.opts)

	if err != nil {
		return 0, 0, errors.Wrap(err, "error while running virtual kubelet")

	}

	return metricsPort, k8sPort, nil
}

// CreateProvider creates the node-cli (virtual kubelet) command
func CreateProvider(env *env.ApateletEnvironment, res *scenario.NodeResources, k8sPort int, metricsPort int, store *store.Store) (*VK, error) {
	op, err := opts.FromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get options from env")
	}

	name := baseName + "-" + res.UUID.String()
	op.KubeConfigPath = env.KubeConfigLocation
	op.ListenPort = 0                      //int32(k8sPort)
	op.MetricsAddr = ":" + strconv.Itoa(0) //metricsPort)
	op.Provider = baseName
	op.NodeName = name

	nodeInfo, err := kubernetes.NewNodeInfo("apatelet", "agent", name, k8sVersion, res.Selector, metricsPort)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes node info")
	}

	providerStore := provider.NewStore()
	providerStore.Register(baseName, func(cfg provider.InitConfig) (provider.Provider, error) {
		cfg.DaemonPort = int32(k8sPort)
		return NewProvider(podmanager.New(), NewStats(), res, cfg, nodeInfo, store), nil
	})

	//node, err := cli.New(ctx,
	//	cli.WithProvider(baseName, func(cfg provider.InitConfig) (provider.Provider, error) {
	//		cfg.DaemonPort = int32(k8sPort)
	//		return NewProvider(podmanager.New(), NewStats(), res, cfg, nodeInfo, store), nil
	//	}),
	//	cli.WithBaseOpts(op),
	//)

	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to create new virtual kubelet provider")
	//}

	return &VK{
		st:   providerStore,
		opts: op,
	}, nil
}

// NewProvider returns the provider but with the vk type instead of our own.
func NewProvider(pods podmanager.PodManager, stats *Stats, resources *scenario.NodeResources, cfg provider.InitConfig, nodeInfo kubernetes.NodeInfo, store *store.Store) provider.Provider {
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
