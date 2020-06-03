package provider

import (
	"context"
	"log"
	"strconv"
	"time"

	nodeconfigv1 "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/apis/nodeconfiguration/v1"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/condition"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/internal/network"
)

const (
	memThresh      = 0.85
	diskThresh     = 0.85
	diskFullThresh = 0.96
	updateInterval = 1 * time.Minute
)

type nodeConditions struct {
	ready          condition.NodeCondition
	outOfDisk      condition.NodeCondition
	memoryPressure condition.NodeCondition
	diskPressure   condition.NodeCondition

	// Unused conditions, may be implement in a later version
	networkUnavailable condition.NodeCondition
	pidPressure        condition.NodeCondition
}

func (p *Provider) getPingResponse() (scenario.Response, error) {
	rawFlag, err := (*p.Store).GetNodeFlag(events.NodePingResponse)
	if err != nil {
		return scenario.ResponseUnset, errors.Errorf("unable to retrieve ping flag %v", err)
	}

	flag, ok := rawFlag.(scenario.Response)
	if !ok {
		return scenario.ResponseUnset, errors.Errorf("invalid ping flag %v", rawFlag)
	}

	return flag, nil
}

// Ping will react to ping based on the given set flag
func (p *Provider) Ping(ctx context.Context) error {
	if err := p.runLatency(ctx); err != nil {
		err = errors.Wrap(err, "failed to run latency (Ping)")
		log.Println(err)
		return err
	}

	flag, err := p.getPingResponse()
	if err != nil {
		return errors.Wrap(err, "getting ping response failed")
	}

	switch flag {
	case scenario.ResponseUnset:
		fallthrough // If unset, act as if it's normal
	case scenario.ResponseNormal:
		return errors.Wrap(ctx.Err(), "context canceled while sending normal response")
	case scenario.ResponseTimeout:
		<-ctx.Done()
		return ctx.Err() // Don't wrap, this context should be closing here
	case scenario.ResponseError:
		return emulationError("ping expected error")
	default:
		return errors.Errorf("invalid response flag %v", flag)
	}
}

// NotifyNodeStatus sets the function we can use to update the status within kubernetes
func (p *Provider) NotifyNodeStatus(ctx context.Context, cb func(*corev1.Node)) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(updateInterval):
				p.updateConditions(ctx, cb)
			}
		}
	}()
}

// ConfigureNode enables a provider to configure the node object that will be used for Kubernetes.
func (p *Provider) ConfigureNode(_ context.Context, node *corev1.Node) {
	node.Spec = p.spec()
	node.ObjectMeta = p.objectMeta()
	node.Status = p.nodeStatus()
	p.Node = node.DeepCopy()
}

func (p *Provider) updateConditions(ctx context.Context, cb func(*corev1.Node)) {
	// First check if the conditions should be updated
	flag, err := p.getPingResponse()
	if err != nil {
		log.Printf("unable to get ping response for updating conditions: %v", err)
		return
	}
	if flag != scenario.ResponseUnset && flag != scenario.ResponseNormal {
		return
	}

	stats, err := p.GetStatsSummary(ctx)
	if err != nil {
		log.Printf("failed to update node conditions: %v", err)
		return
	}

	// Set bools
	memPressure := float32(*stats.Node.Memory.UsageBytes) > float32(p.Resources.Memory)*memThresh
	diskPressure := float32(*stats.Node.Fs.UsedBytes) > float32(p.Resources.Storage)*diskThresh
	diskFull := float32(*stats.Node.Fs.UsedBytes) > float32(p.Resources.Storage)*diskFullThresh

	// Set conditions and update node
	p.Node.Status.Conditions = []corev1.NodeCondition{
		p.Conditions.ready.Update(!diskFull),
		p.Conditions.outOfDisk.Update(diskFull),
		p.Conditions.memoryPressure.Update(memPressure),
		p.Conditions.diskPressure.Update(diskPressure),
		p.Conditions.networkUnavailable.Update(false),
		p.Conditions.pidPressure.Update(false),
	}
	cb(p.Node)
}

func (p *Provider) nodeConditions() []corev1.NodeCondition {
	// Return static thumbs-up values for all conditions.
	return []corev1.NodeCondition{
		p.Conditions.ready.Get(),
		p.Conditions.outOfDisk.Get(),
		p.Conditions.memoryPressure.Get(),
		p.Conditions.diskPressure.Get(),
		p.Conditions.networkUnavailable.Get(),
		p.Conditions.pidPressure.Get(),
	}
}

func (p *Provider) nodeStatus() corev1.NodeStatus {
	return corev1.NodeStatus{
		NodeInfo: corev1.NodeSystemInfo{
			Architecture:   "amd64",
			KubeletVersion: p.NodeInfo.Version,
		},
		DaemonEndpoints: p.nodeDaemonEndpoints(),
		Addresses:       p.addresses(),
		Capacity:        p.capacity(),
		Conditions:      p.nodeConditions(),
	}
}

func (p *Provider) objectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: p.NodeInfo.Name,
		Labels: map[string]string{
			"type":                     p.NodeInfo.NodeType,
			"kubernetes.io/role":       p.NodeInfo.Role,
			"kubernetes.io/hostname":   p.NodeInfo.Name,
			"metrics_port":             strconv.Itoa(p.NodeInfo.MetricsPort),
			nodeconfigv1.EmulatedLabel: nodeconfigv1.EmulatedLabelValue,

			nodeconfigv1.NodeConfigurationLabelNamespace: p.NodeInfo.Namespace,
			nodeconfigv1.NodeConfigurationLabel:          p.NodeInfo.Selector,
		},
	}
}

func (p *Provider) spec() corev1.NodeSpec {
	if p.DisableTaints {
		return corev1.NodeSpec{}
	}

	return corev1.NodeSpec{
		Taints: []corev1.Taint{
			{
				Key:    nodeconfigv1.EmulatedLabel,
				Effect: corev1.TaintEffectNoSchedule,
			},
		},
	}
}

func (p *Provider) addresses() []corev1.NodeAddress {
	externalAddress, err := network.GetExternalAddress()
	if err != nil {
		log.Printf("error while retrieving ip addresses for node: %v\n", err)
		return []corev1.NodeAddress{}
	}

	return []corev1.NodeAddress{
		{
			Type:    "InternalIP",
			Address: externalAddress,
		},
		{
			Type:    "ExternalIP",
			Address: externalAddress,
		},
	}
}

func (p *Provider) nodeDaemonEndpoints() corev1.NodeDaemonEndpoints {
	return corev1.NodeDaemonEndpoints{
		KubeletEndpoint: corev1.DaemonEndpoint{
			Port: p.Cfg.DaemonPort,
		},
	}
}

func (p *Provider) capacity() corev1.ResourceList {
	var cpu resource.Quantity
	cpu.Set(p.Resources.CPU)

	var mem resource.Quantity
	mem.Set(p.Resources.Memory)

	var pods resource.Quantity
	pods.Set(p.Resources.MaxPods)

	var storage resource.Quantity
	storage.Set(p.Resources.Storage)

	var ephemeralStorage resource.Quantity
	ephemeralStorage.Set(p.Resources.EphemeralStorage)

	return corev1.ResourceList{
		corev1.ResourceCPU:              cpu,
		corev1.ResourceMemory:           mem,
		corev1.ResourcePods:             pods,
		corev1.ResourceStorage:          storage,
		corev1.ResourceEphemeralStorage: ephemeralStorage,
	}
}
