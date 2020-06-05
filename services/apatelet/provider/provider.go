// Package provider handles the interaction with the virtual kubelet library (and node-cli)
package provider

import (
	"context"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/node"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	root "github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/commands"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/opts"
	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/node-cli/provider"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/condition"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
	baseName   = "apatelet"
)

// Provider implements the node-cli (virtual kubelet) interface for a virtual kubelet provider
type Provider struct {
	Pods        podmanager.PodManager // the pods currently used
	Store       *store.Store          // the apatelet store
	Environment env.ApateletEnvironment

	Cfg           *provider.InitConfig // the initial provider config
	DisableTaints bool                 // whether to disable taints

	Stats *Stats // statistics contain static statistics

	Node      *corev1.Node            // the reference to "ourselves"
	NodeInfo  *node.Info              // static node information sent to kubernetes
	Resources *scenario.NodeResources // static resource information sent to kubernetes

	Conditions nodeConditions // a wrapper around kubernetes conditions
}

// VirtualKubelet is a struct containing everything needed to start virtual kubelet
type VirtualKubelet struct {
	st   *provider.Store
	opts *opts.Opts

	info *node.Info
}

//nolint as lint does not recognise the first context is indeed the correct context
// Run starts the virtual kubelet
func (vk *VirtualKubelet) Run(ctx context.Context, originalCtx context.Context) (int, int, error) {
	metricsPort, k8sPort, err := root.RunRootCommand(originalCtx, ctx, vk.st, vk.opts)

	if err != nil {
		return 0, 0, errors.Wrap(err, "error while running virtual kubelet")
	}

	// Update metrics port
	vk.info.MetricsPort = metricsPort

	return metricsPort, k8sPort, nil
}

// CreateProvider creates the node-cli (virtual kubelet) command
func CreateProvider(env *env.ApateletEnvironment, res *scenario.NodeResources, store *store.Store) (*VirtualKubelet, error) {
	op, err := opts.FromEnv()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get options from env")
	}

	name := baseName + "-" + res.UUID.String()
	op.KubeConfigPath = env.KubeConfigLocation
	op.ListenPort = int32(0)
	op.MetricsAddr = ":0"
	op.Provider = baseName
	op.NodeName = name

	nodeInfo, err := node.NewInfo("apatelet", "agent", name, k8sVersion, res.Label)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes node info")
	}

	providerStore := provider.NewStore()
	providerStore.Register(baseName, func(cfg *provider.InitConfig) (provider.Provider, error) {
		return NewProvider(podmanager.New(), NewStats(), res, cfg, &nodeInfo, store, env.DisableTaints, *env), nil
	})

	return &VirtualKubelet{
		st:   providerStore,
		opts: op,
		info: &nodeInfo,
	}, nil
}

// NewProvider returns the provider but with the vk type instead of our own.
func NewProvider(pods podmanager.PodManager, nodeStats *Stats, resources *scenario.NodeResources, cfg *provider.InitConfig, nodeInfo *node.Info, store *store.Store, disableTaints bool, environment env.ApateletEnvironment) provider.Provider {
	p := &Provider{
		Pods:        pods,
		Store:       store,
		Environment: environment,

		Cfg:           cfg,
		DisableTaints: disableTaints,

		Stats: nodeStats,

		NodeInfo:  nodeInfo,
		Resources: resources,

		Conditions: nodeConditions{
			ready:              condition.New(true, corev1.NodeReady),
			outOfDisk:          condition.New(false, corev1.NodeOutOfDisk),
			memoryPressure:     condition.New(false, corev1.NodeMemoryPressure),
			diskPressure:       condition.New(false, corev1.NodeDiskPressure),
			networkUnavailable: condition.New(false, corev1.NodeNetworkUnavailable),
			pidPressure:        condition.New(false, corev1.NodePIDPressure),
		},
	}

	(*store).AddPodListener(events.PodResources, func(obj interface{}) {
		p.updateStatsSummary()
	})

	p.updateStatsSummary()

	return p
}
