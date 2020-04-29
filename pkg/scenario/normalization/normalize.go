// Package normalization provides functions to normalize and decode public scenarios.
package normalization

import (
	"time"

	"github.com/google/uuid"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/controlplane/events"
	apiScenario "github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"github.com/docker/go-units"
)

// NormalizeScenario takes a public scenario and turns it into a private scenario.
// Normalizes the structure and resolves named references.
func NormalizeScenario(scenario *controlplane.PublicScenario) (*apatelet.ApateletScenario, []NodeResources, error) {
	r := apatelet.ApateletScenario{}

	nodeResources := make([]NodeResources, 0)
	uuidsPerNodeGroup := make(map[string][]uuid.UUID)

	// First make a lookup mapping nodeType strings to node types.
	// This makes later lookup O(1)
	nodeTypeName := make(map[string]*controlplane.Node)
	for _, nodeType := range scenario.GetNodes() {
		nodeTypeName[nodeType.NodeType] = nodeType
	}

	for _, nodeGroup := range scenario.NodeGroups {
		for i := 0; i < int(nodeGroup.Amount); i++ {
			id := uuid.New()

			nodeType := nodeTypeName[nodeGroup.NodeType]
			memory, err := units.RAMInBytes(nodeType.Memory)
			if err != nil {
				return nil, nil, err
			}

			nodeResources = append(nodeResources, NodeResources{
				id,
				memory,
				nodeType.Cpu,
				nodeType.MaxPods,
			})

			uuidsPerNodeGroup[nodeGroup.GroupName] = append(uuidsPerNodeGroup[nodeGroup.GroupName], id)
		}
	}

	var tasks []*apatelet.Task

	for _, task := range scenario.Tasks {
		timestamp, err := desugarTimestamp(task.Time)
		if err != nil {
			return nil, nil, err
		}

		// Decode the "all" node name, also verify that all names in the nodeSet exist and
		// that there are no duplicates in the set.
		nodeGroupNames, err := desugarNodeGroups(task.NodeGroups, scenario.NodeGroups)
		if err != nil {
			return nil, nil, err
		}

		var nodeSet []string

		for _, name := range nodeGroupNames {
			for _, nodeUUID := range uuidsPerNodeGroup[name] {
				nodeSet = append(nodeSet, nodeUUID.String())
			}
		}

		newTask := &apatelet.Task{
			Name:       task.Name,
			RevertTask: task.Revert,
			Timestamp:  int32(timestamp),
			NodeSet:    nodeSet,
		}

		translateEvent(task, newTask)

		tasks = append(tasks, newTask)
	}

	r.Task = tasks

	return &r, nodeResources, nil
}

func translateEvent(originalTask *controlplane.Task, newTask *apatelet.Task) {
	switch x := originalTask.Event.(type) {
	// Node events
	case *controlplane.Task_NodeFailure:
		nodeEvent := &apatelet.Task_NodeEvent{}
		nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
			LifecycleState: &apatelet.LifecycleState{
				CreateAction:                 apiScenario.LifecycleAction_TIMEOUT,
				CreateActionPercentage:       100,
				UpdateAction:                 apiScenario.LifecycleAction_TIMEOUT,
				UpdateActionPercentage:       100,
				DeleteAction:                 apiScenario.LifecycleAction_TIMEOUT,
				DeleteActionPercentage:       100,
				GetPodAction:                 apiScenario.LifecycleAction_TIMEOUT,
				GetPodActionPercentage:       100,
				GetPodStatusAction:           apiScenario.LifecycleAction_TIMEOUT,
				GetPodStatusActionPercentage: 100,
			},
			GetPodsAction:           apiScenario.LifecycleAction_TIMEOUT,
			GetPodsActionPercentage: 100,
			PingAction:              apiScenario.LifecycleAction_TIMEOUT,
			PingActionPercentage:    100,
		}
		newTask.Event = nodeEvent

	case *controlplane.Task_NetworkLatency:
		nodeEvent := &apatelet.Task_NodeEvent{}
		nodeEvent.NodeEvent.NodeState.AddedLatencyState = &apatelet.NodeState_AddedLatencyState{
			AddedLatencyEnabled: true,
			AddedLatencyMsec:    x.NetworkLatency.GetLatencyMsec(),
		}
		newTask.Event = nodeEvent

	case *controlplane.Task_TimeoutKeepHeartbeat:
		nodeEvent := &apatelet.Task_NodeEvent{}
		nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
			LifecycleState: &apatelet.LifecycleState{
				CreateAction:                 apiScenario.LifecycleAction_TIMEOUT,
				CreateActionPercentage:       100,
				UpdateAction:                 apiScenario.LifecycleAction_TIMEOUT,
				UpdateActionPercentage:       100,
				DeleteAction:                 apiScenario.LifecycleAction_TIMEOUT,
				DeleteActionPercentage:       100,
				GetPodAction:                 apiScenario.LifecycleAction_TIMEOUT,
				GetPodActionPercentage:       100,
				GetPodStatusAction:           apiScenario.LifecycleAction_TIMEOUT,
				GetPodStatusActionPercentage: 100,
			},
			GetPodsAction:           apiScenario.LifecycleAction_TIMEOUT,
			GetPodsActionPercentage: 100,
		}
		newTask.Event = nodeEvent

	case *controlplane.Task_NoTimeoutNoHeartbeat:
		nodeEvent := &apatelet.Task_NodeEvent{}
		nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
			PingAction:           apiScenario.LifecycleAction_TIMEOUT,
			PingActionPercentage: 100,
		}
		newTask.Event = nodeEvent

	case *controlplane.Task_NodeLifecycleState:
		nodeEvent := &apatelet.Task_NodeEvent{}

		switch x.NodeLifecycleState.Type {
		case events.LifecycleType_CREATE_POD:
			nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					CreateAction:           x.NodeLifecycleState.Action,
					CreateActionPercentage: x.NodeLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_UPDATE_POD:
			nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					UpdateAction:           x.NodeLifecycleState.Action,
					UpdateActionPercentage: x.NodeLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_DELETE_POD:
			nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					DeleteAction:           x.NodeLifecycleState.Action,
					DeleteActionPercentage: x.NodeLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_GET_POD:
			nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					GetPodAction:           x.NodeLifecycleState.Action,
					GetPodActionPercentage: x.NodeLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_GET_POD_STATUS:
			nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					GetPodStatusAction:           x.NodeLifecycleState.Action,
					GetPodStatusActionPercentage: x.NodeLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_GET_PODS:
			nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
				GetPodsAction:           x.NodeLifecycleState.Action,
				GetPodsActionPercentage: x.NodeLifecycleState.Percentage,
			}

		case events.LifecycleType_PING:
			nodeEvent.NodeEvent.NodeState.NodeLifecycleState = &apatelet.NodeState_NodeLifecycleState{
				PingAction:           x.NodeLifecycleState.Action,
				PingActionPercentage: x.NodeLifecycleState.Percentage,
			}
		}
		newTask.Event = nodeEvent

	case *controlplane.Task_ResourcePressure:
		nodeEvent := &apatelet.Task_NodeEvent{}
		nodeEvent.NodeEvent.NodeState.ResourceState = &apatelet.NodeState_ResourceState{
			EnableResourceAlteration: true,
			CpuUsage:                 x.ResourcePressure.GetCpuUsage(),
			MemoryUsage:              x.ResourcePressure.GetMemoryUsage(),
			StorageUsage:             x.ResourcePressure.GetStorageUsage(),
			EphemeralStorageUsage:    x.ResourcePressure.GetEphemeralStorageUsage(),
		}
		newTask.Event = nodeEvent

	// Pod events
	case *controlplane.Task_PodLifecycleState:
		podEvent := &apatelet.Task_PodEvent{}

		switch x.PodLifecycleState.Type {
		case events.LifecycleType_CREATE_POD:
			podEvent.PodEvent.PodState.PodLifecycleState = &apatelet.PodState_PodLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					CreateAction:           x.PodLifecycleState.Action,
					CreateActionPercentage: x.PodLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_UPDATE_POD:
			podEvent.PodEvent.PodState.PodLifecycleState = &apatelet.PodState_PodLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					UpdateAction:           x.PodLifecycleState.Action,
					UpdateActionPercentage: x.PodLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_DELETE_POD:
			podEvent.PodEvent.PodState.PodLifecycleState = &apatelet.PodState_PodLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					DeleteAction:           x.PodLifecycleState.Action,
					DeleteActionPercentage: x.PodLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_GET_POD:
			podEvent.PodEvent.PodState.PodLifecycleState = &apatelet.PodState_PodLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					GetPodAction:           x.PodLifecycleState.Action,
					GetPodActionPercentage: x.PodLifecycleState.Percentage,
				},
			}

		case events.LifecycleType_GET_POD_STATUS:
			podEvent.PodEvent.PodState.PodLifecycleState = &apatelet.PodState_PodLifecycleState{
				LifecycleState: &apatelet.LifecycleState{
					GetPodStatusAction:           x.PodLifecycleState.Action,
					GetPodStatusActionPercentage: x.PodLifecycleState.Percentage,
				},
			}
		}
		newTask.Event = podEvent

	case *controlplane.Task_PodStatusUpdate:
		podEvent := &apatelet.Task_PodEvent{}
		podEvent.PodEvent.PodState.PodStatus = x.PodStatusUpdate.NewStatus
		newTask.Event = podEvent

	case *controlplane.Task_PodStartTimeUpdate:
		podEvent := &apatelet.Task_PodEvent{}
		podEvent.PodEvent.PodState.StartTime = x.PodStartTimeUpdate.NewStartTime
		newTask.Event = podEvent
	}
}

func desugarTimestamp(t string) (int, error) {
	duration, err := time.ParseDuration(t)
	if err != nil {
		return 0, err
	}

	return int(duration.Milliseconds()), nil
}
