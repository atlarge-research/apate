package provider

import (
	"context"
	"errors"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

type responseArgs struct {
	ctx      context.Context
	provider *Provider
	action   func() (interface{}, error)
}

func podResponse(args responseArgs, label string, responseFlag events.PodEventFlag) (interface{}, bool, error) {
	iflag, err := (*args.provider.store).GetPodFlag(label, responseFlag)
	if err != nil {
		return nil, false, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, false, errors.New("podResponse couldn't cast flag to response")
	}

	switch flag {
	case scenario.Response_RESPONSE_NORMAL:
		res, err := args.action()
		return res, true, err
	case scenario.Response_RESPONSE_TIMEOUT:
		<-args.ctx.Done()
		return nil, true, nil
	case scenario.Response_RESPONSE_ERROR:
		return nil, true, errors.New("podResponse expected error")
	case scenario.Response_RESPONSE_UNSET:
		return nil, false, nil
	default:
		return nil, false, errors.New("podResponse invalid scenario")
	}
}

// nodeResponse checks the passed flags and calls the passed function on success
func nodeResponse(args responseArgs, responseFlag events.NodeEventFlag) (interface{}, error) {
	iflag, err := (*args.provider.store).GetNodeFlag(responseFlag)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, errors.New("nodeResponse couldn't cast flag to response")
	}

	switch flag {
	case scenario.Response_RESPONSE_UNSET:
		fallthrough // If unset, act as if it's normal
	case scenario.Response_RESPONSE_NORMAL:
		res, err := args.action()
		return res, err
	case scenario.Response_RESPONSE_TIMEOUT:
		<-args.ctx.Done()
		return nil, nil
	case scenario.Response_RESPONSE_ERROR:
		return nil, errors.New("nodeResponse expected error")
	default:
		return nil, errors.New("nodeResponse invalid scenario")
	}
}

// podAndNodeResponse first calls podResponse and if pod response is not set calls nodeResponse
func podAndNodeResponse(args responseArgs, podLabel string, podResponseFlag events.PodEventFlag, nodeResponseFlag events.NodeEventFlag) (interface{}, error) {
	pod, performedAction, err := podResponse(args, podLabel, podResponseFlag)

	if err != nil {
		return nil, err
	}

	if performedAction {
		return pod, nil
	}

	return nodeResponse(args, nodeResponseFlag)
}
