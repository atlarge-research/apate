package provider

import (
	"context"
	corev1 "k8s.io/api/core/v1"

	"github.com/pkg/errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"

	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

// An emulation error is an error occurring because the system is emulating an error. This type of error is completely expected.
// This error should never be wrapped
type emulationError string

func (e emulationError) Error() string {
	return string(e)
}

func (e emulationError) Expected() bool {
	return true
}

type responseArgs struct {
	ctx      context.Context
	provider *Provider
	action   func() (interface{}, error)
}

func podResponse(args responseArgs, label string, pod *corev1.Pod, responseFlag events.PodEventFlag) (interface{}, bool, error) {
	iflag, err := (*args.provider.Store).GetPodFlag(label, pod, responseFlag)
	if err != nil {
		return nil, false, errors.Wrap(err, "failed to get pod flag")
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, false, errors.Errorf("couldn't cast %v to response", flag)
	}

	switch flag {
	case scenario.ResponseNormal:
		res, err := args.action()
		return res, true, errors.Wrap(err, "failed to execute pod response action")
	case scenario.ResponseTimeout:
		<-args.ctx.Done()
		return nil, true, nil
	case scenario.ResponseError:
		return nil, true, emulationError("podResponse expected error")
	case scenario.ResponseUnset:
		return nil, false, nil
	default:
		return nil, false, errors.New("podResponse invalid scenario")
	}
}

// nodeResponse checks the passed flags and calls the passed function on success
func nodeResponse(args responseArgs, responseFlag events.NodeEventFlag) (interface{}, error) {
	iflag, err := (*args.provider.Store).GetNodeFlag(responseFlag)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get node flag")
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, errors.Errorf("couldn't cast %v to response", flag)
	}

	switch flag {
	case scenario.ResponseUnset:
		fallthrough // If unset, act as if it's normal
	case scenario.ResponseNormal:
		res, err := args.action()
		return res, errors.Wrap(err, "failed to execute node response action")
	case scenario.ResponseTimeout:
		<-args.ctx.Done()
		return nil, nil
	case scenario.ResponseError:
		return nil, emulationError("nodeResponse expected error")
	default:
		return nil, errors.New("nodeResponse invalid scenario")
	}
}

// podAndNodeResponse first calls podResponse and if pod response is not set calls nodeResponse
func podAndNodeResponse(args responseArgs, podLabel string, pod *corev1.Pod, podResponseFlag events.PodEventFlag, nodeResponseFlag events.NodeEventFlag) (interface{}, error) {
	podFlag, performedAction, err := podResponse(args, podLabel, pod, podResponseFlag)

	if err != nil {
		if IsExpected(err) {
			return nil, err
		}

		return nil, errors.Wrap(err, "failed during pod response (not going to try node response)")
	}

	if performedAction {
		return podFlag, nil
	}

	node, err := nodeResponse(args, nodeResponseFlag)
	if err != nil {
		if IsExpected(err) {
			return nil, err
		}

		return nil, errors.Wrap(err, "failed during pod response")
	}
	return node, nil
}
