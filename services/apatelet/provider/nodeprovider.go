package provider

import (
	"context"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/network"
)

// Ping TODO
func (p *Provider) Ping(ctx context.Context) error {
	return ctx.Err()
}

// NotifyNodeStatus TODO
func (p *Provider) NotifyNodeStatus(_ context.Context, _ func(*corev1.Node)) {
	// TODO
}

// ConfigureNode enables a provider to configure the node object that will be used for Kubernetes.
func (p *Provider) ConfigureNode(_ context.Context, node *corev1.Node) {
	node.Spec = p.spec()
	node.ObjectMeta = p.objectMeta()
	node.Status = p.nodeStatus()
}

func (p *Provider) nodeConditions() []corev1.NodeCondition {
	lastHeartbeatTime := metav1.Now()
	lastTransitionTime := metav1.Now()
	lastTransitionReason := "Apatelet is ready"
	lastTransitionMessage := "ok"

	// Return static thumbs-up values for all conditions.
	return []corev1.NodeCondition{
		{
			Type:               corev1.NodeReady,
			Status:             corev1.ConditionTrue,
			LastHeartbeatTime:  lastHeartbeatTime,
			LastTransitionTime: lastTransitionTime,
			Reason:             lastTransitionReason,
			Message:            lastTransitionMessage,
		},
		{
			Type:               corev1.NodeOutOfDisk,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  lastHeartbeatTime,
			LastTransitionTime: lastTransitionTime,
			Reason:             lastTransitionReason,
			Message:            lastTransitionMessage,
		},
		{
			Type:               corev1.NodeMemoryPressure,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  lastHeartbeatTime,
			LastTransitionTime: lastTransitionTime,
			Reason:             lastTransitionReason,
			Message:            lastTransitionMessage,
		},
		{
			Type:               corev1.NodeDiskPressure,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  lastHeartbeatTime,
			LastTransitionTime: lastTransitionTime,
			Reason:             lastTransitionReason,
			Message:            lastTransitionMessage,
		},
		{
			Type:               corev1.NodeNetworkUnavailable,
			Status:             corev1.ConditionFalse,
			LastHeartbeatTime:  lastHeartbeatTime,
			LastTransitionTime: lastTransitionTime,
			Reason:             lastTransitionReason,
			Message:            lastTransitionMessage,
		},
		{
			Type:               "KubeletConfigOk",
			Status:             corev1.ConditionTrue,
			LastHeartbeatTime:  lastHeartbeatTime,
			LastTransitionTime: lastTransitionTime,
			Reason:             lastTransitionReason,
			Message:            lastTransitionMessage,
		},
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
