// Package provider handles the interaction with the virtual kubelet library (and node-cli)
package provider

import (
	"context"
	"strconv"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/env"

	corev1 "k8s.io/api/core/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/condition"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/podmanager"

	cli "github.com/virtual-kubelet/node-cli"
	"github.com/virtual-kubelet/node-cli/opts"
	"github.com/virtual-kubelet/node-cli/provider"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

var (
	k8sVersion = "v1.15.2" // This should follow the version of k8s.io/kubernetes we are importing
	baseName   = "apatelet"
)

// Provider implements the node-cli (virtual kubelet) interface for a virtual kubelet provider
type Provider struct {
	Pods  podmanager.PodManager // the pods currently used
	Store *store.Store          // the apatelet store

	Cfg           provider.InitConfig // the initial provider config
	DisableTaints bool                // whether to disable taints

	Stats *Stats // statistics contain static statistics

	Node      *corev1.Node            // the reference to "ourselves"
	NodeInfo  kubernetes.NodeInfo     // static node information sent to kubernetes
	Resources *scenario.NodeResources // static resource information sent to kubernetes

	Conditions nodeConditions // a wrapper around kubernetes conditions
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

	nodeInfo, err := kubernetes.NewNodeInfo("apatelet", "agent", name, k8sVersion, res.Label, metricsPort)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create kubernetes node info")
	}

	node, err := cli.New(ctx,
		cli.WithProvider(baseName, func(cfg provider.InitConfig) (provider.Provider, error) {
			cfg.DaemonPort = int32(k8sPort)
			return NewProvider(podmanager.New(), NewStats(), res, cfg, nodeInfo, store, env.DisableTaints), nil
		}),
		cli.WithBaseOpts(op),
	)

	if err != nil {
		return nil, errors.Wrap(err, "failed to create new virtual kubelet provider")
	}

	return node, nil
}

// NewProvider returns the provider but with the vk type instead of our own.
func NewProvider(pods podmanager.PodManager, nodeStats *Stats, resources *scenario.NodeResources, cfg provider.InitConfig, nodeInfo kubernetes.NodeInfo, store *store.Store, disableTaints bool) provider.Provider {
	p := &Provider{
		Pods:  pods,
		Store: store,

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
