package provider

import (
	"context"
	corev1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

type responseArgs struct {
	ctx      context.Context
	provider *Provider
	action   func() (interface{}, error)
}

func getPodResponseFlag(args responseArgs, label string, podEventFlag events.PodEventFlag) (scenario.Response, error) {
	iflag, err := (*args.provider.Store).GetPodFlag(label, podEventFlag)
	if err != nil {
		return scenario.ResponseUnset, errors.Wrapf(err, "failed to get pod flag %v", podEventFlag)
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return scenario.ResponseUnset, errors.Errorf("couldn't cast %v to response", flag)
	}

	return flag, nil
}

func getNodeResponseFlag(args responseArgs, nodeEventFlag events.NodeEventFlag) (scenario.Response, error) {
	iflag, err := (*args.provider.Store).GetNodeFlag(nodeEventFlag)
	if err != nil {
		return scenario.ResponseUnset, errors.Wrapf(err, "failed to get node flag %v", nodeEventFlag)
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return scenario.ResponseUnset, errors.Errorf("couldn't cast %v to response", flag)
	}

	return flag, nil
}

// nodeResponse checks the passed flags and calls the passed function on success
func nodeResponse(args responseArgs, nodeEventFlag events.NodeEventFlag) (interface{}, error) {
	flag, err := getNodeResponseFlag(args, nodeEventFlag)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve node flag")
	}

	switch flag {
	case scenario.ResponseTimeout:
		<-args.ctx.Done()
		return nil, nil
	case scenario.ResponseError:
		return nil, emulationError("nodeResponse expected error")
	case scenario.ResponseUnset:
		fallthrough
	case scenario.ResponseNormal:
		fallthrough
	default:
		res, err := args.action()
		return res, errors.Wrap(err, "failed to execute node response action")
	}
}

func getCorrespondingNodeEventFlag(podEventFlag events.PodEventFlag) (events.NodeEventFlag, error) {
	switch podEventFlag {
	case events.PodCreatePodResponse:
		return events.NodeCreatePodResponse, nil
	case events.PodUpdatePodResponse:
		return events.NodeUpdatePodResponse, nil
	case events.PodDeletePodResponse:
		return events.NodeDeletePodResponse, nil
	case events.PodGetPodResponse:
		return events.NodeGetPodResponse, nil
	case events.PodGetPodStatusResponse:
		return events.NodeGetPodStatusResponse, nil
	default:
		return -1, errors.Errorf("pod event flag %v cannot be translated into node flag", podEventFlag)
	}
}

// podResponse determines how a pod action should respond, also based on flags on node level
func podResponse(args responseArgs, podLabel string, podEventFlag events.PodEventFlag) (interface{}, error) {
	podFlag, err := getPodResponseFlag(args, podLabel, podEventFlag)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve pod flag")
	}

	nodeEventFlag, err := getCorrespondingNodeEventFlag(podEventFlag)
	if err != nil {
		return nil, errors.Wrap(err, "error translating pod event flag to node event flag")
	}

	nodeFlag, err := getNodeResponseFlag(args, nodeEventFlag)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve node flag")
	}

	switch {
	case podFlag == scenario.ResponseTimeout || nodeFlag == scenario.ResponseTimeout:
		// first try for timeout
		<-args.ctx.Done()
		return nil, nil
	case podFlag == scenario.ResponseError || nodeFlag == scenario.ResponseError:
		// then error
		return nil, emulationError("nodeResponse expected error")
	default:
		// then normal. Unset is also normal, as is any invalid type (which shouldn't even be possible)
		res, err := args.action()
		return res, errors.Wrap(err, "failed to execute node response action")
	}
}
