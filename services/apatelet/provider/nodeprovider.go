package provider

import (
	"context"
	"log"
	"strconv"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/provider/condition"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/network"
)

const (
	memThresh      = 0.85
	diskThresh     = 0.85
	diskFullThresh = 0.96
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
	rawFlag, err := (*p.store).GetNodeFlag(events.NodePingResponse)
	if err != nil {
		return scenario.ResponseUnset, errors.Errorf("unable to retrieve ping flag: %v", err)
	}

	flag, ok := rawFlag.(scenario.Response)
	if !ok {
		return scenario.ResponseUnset, errors.Errorf("invalid ping flag: %v", rawFlag)
	}

	return flag, nil
}

// Ping will react to ping based on the given set flag
func (p *Provider) Ping(ctx context.Context) error {
	flag, err := p.getPingResponse()
	if err != nil {
		return errors.Wrap(err, "getting ping response failed")
	}

	switch flag {
	case scenario.ResponseUnset:
		fallthrough // If unset, act as if it's normal
	case scenario.ResponseNormal:
		return ctx.Err()
	case scenario.ResponseTimeout:
		<-ctx.Done()
		return ctx.Err()
	case scenario.ResponseError:
		return errors.Errorf("ping expected error")
	default:
		return errors.Errorf("invalid response flag: %v", flag)
	}
}

// NotifyNodeStatus sets the function we can use to update the status within kubernetes
func (p *Provider) NotifyNodeStatus(_ context.Context, cb func(*corev1.Node)) {
	p.updateStatus = cb
}

// ConfigureNode enables a provider to configure the node object that will be used for Kubernetes.
func (p *Provider) ConfigureNode(_ context.Context, node *corev1.Node) {
	node.Spec = p.spec()
	node.ObjectMeta = p.objectMeta()
	node.Status = p.nodeStatus()
	p.node = node.DeepCopy()
}

func (p *Provider) updateConditions(ctx context.Context) {
	// First check if the conditions should be updated
	flag, err := p.getPingResponse()
	if err != nil {
		log.Printf("unable to get ping response for updating conditions: %v", err)
		return
	}
	if flag != scenario.ResponseUnset && flag != scenario.ResponseNormal {
		return //TODO: Should we log this? Might result in some spam logging..
	}

	stats, err := p.GetStatsSummary(ctx)
	if err != nil {
		// TODO: What to do now?
		log.Printf("failed to update node conditions: %v", err)
		return
	}

	// Set bools
	memPressure := float64(*stats.Node.Memory.UsageBytes) > float64(p.resources.Memory)*memThresh
	diskPressure := float64(*stats.Node.Fs.UsedBytes) > float64(p.resources.Storage)*diskThresh
	diskFull := float64(*stats.Node.Fs.UsedBytes) > float64(p.resources.Storage)*diskFullThresh

	// Set conditions and update node
	p.node.Status.Conditions = []corev1.NodeCondition{
		p.conditions.ready.Update(!diskFull),
		p.conditions.outOfDisk.Update(diskFull),
		p.conditions.memoryPressure.Update(memPressure),
		p.conditions.diskPressure.Update(diskPressure),
		p.conditions.networkUnavailable.Update(false),
		p.conditions.pidPressure.Update(false),
	}
	p.updateStatus(p.node)
}

func (p *Provider) nodeConditions() []corev1.NodeCondition {
	// Return static thumbs-up values for all conditions.
	return []corev1.NodeCondition{
		p.conditions.ready.Get(),
		p.conditions.outOfDisk.Get(),
		p.conditions.memoryPressure.Get(),
		p.conditions.diskPressure.Get(),
		p.conditions.networkUnavailable.Get(),
		p.conditions.pidPressure.Get(),
	}
}

func (p *Provider) nodeStatus() corev1.NodeStatus {
	return corev1.NodeStatus{
		NodeInfo: corev1.NodeSystemInfo{
			Architecture:   "amd64",
			KubeletVersion: p.nodeInfo.Version,
		},
		DaemonEndpoints: p.nodeDaemonEndpoints(),
		Addresses:       p.addresses(),
		Capacity:        p.capacity(),
		Conditions:      p.nodeConditions(),
	}
}

func (p *Provider) objectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name: p.nodeInfo.Name,
		Labels: map[string]string{
			"type":                   p.nodeInfo.NodeType,
			"kubernetes.io/role":     p.nodeInfo.Role,
			"kubernetes.io/hostname": p.nodeInfo.Name,
			"metrics_port":           strconv.Itoa(p.nodeInfo.MetricsPort),
			"apate":                  p.nodeInfo.Selector,
		},
	}
}

func (p *Provider) spec() corev1.NodeSpec {
	taints := make([]corev1.Taint, 0)
	return corev1.NodeSpec{
		Taints: taints,
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
			Port: p.cfg.DaemonPort,
		},
	}
}

func (p *Provider) capacity() corev1.ResourceList {
	var cpu resource.Quantity
	cpu.Set(p.resources.CPU)

	var mem resource.Quantity
	mem.Set(p.resources.Memory)

	var pods resource.Quantity
	pods.Set(p.resources.MaxPods)

	var storage resource.Quantity
	storage.Set(p.resources.Storage)

	var ephemeralStorage resource.Quantity
	ephemeralStorage.Set(p.resources.EphemeralStorage)

	return corev1.ResourceList{
		corev1.ResourceCPU:              cpu,
		corev1.ResourceMemory:           mem,
		corev1.ResourcePods:             pods,
		corev1.ResourceStorage:          storage,
		corev1.ResourceEphemeralStorage: ephemeralStorage,
	}
}
