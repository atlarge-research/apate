package provider

import (
	"context"
	"math/rand"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/throw"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

const flagNotSetError = store.FlagNotFoundError
const expectedError = throw.Exception("expected error")
const invalidResponse = throw.Exception("invalid response")
const invalidPercentage = throw.Exception("invalid percentage type")
const invalidFlag = throw.Exception("invalid flag type")

type responseArgs struct {
	ctx      context.Context
	provider *Provider
	action   func() (interface{}, error)
}

type podResponseArgs struct {
	name string

	podResponseFlag   events.PodEventFlag
	podPercentageFlag events.PodEventFlag
}

type nodeResponseArgs struct {
	nodeResponseFlag   events.NodeEventFlag
	nodePercentageFlag events.NodeEventFlag
}

type podNodeResponse struct {
	responseArgs
	podResponseArgs
	nodeResponseArgs
}

/**
This file contains helper functions for the provider in the case that the provider has the structure of

Two pod flags:
	* podResponseFlag
	* podResponsePercentageFlag

and two node flags:
	* nodeResponseFlag
	* nodeResponsePercentageFlag

the the helper function will get the flags from the store, check if they are valid,
calculate the percentage and calls the action callback on success.
*/

// podResponse checks the passed flags and calls the passed function on success
func podResponse(args responseArgs, podA podResponseArgs) (interface{}, bool, error) {
	iflag, err := (*args.provider.store).GetPodFlag(podA.name, podA.podResponseFlag)
	if err != nil {
		return nil, false, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, false, invalidFlag
	}

	iflagpercent, err := (*args.provider.store).GetPodFlag(podA.name, podA.podPercentageFlag)
	if err != nil {
		return nil, false, err
	}

	flagpercent, ok := iflagpercent.(int32)
	if !ok {
		return nil, false, invalidPercentage
	}

	if flagpercent < rand.Int31n(int32(100)) {
		res, err := args.action()
		return res, true, err
	}

	switch flag {
	case scenario.Response_NORMAL:
		res, err := args.action()
		return res, true, err
	case scenario.Response_TIMEOUT:
		<-args.ctx.Done()
		return nil, true, nil
	case scenario.Response_ERROR:
		return nil, true, expectedError
	default:
		return nil, false, invalidResponse
	}
}

// nodeResponse checks the passed flags and calls the passed function on success
func nodeResponse(args responseArgs, nodeA nodeResponseArgs) (interface{}, bool, error) {
	iflag, err := (*args.provider.store).GetNodeFlag(nodeA.nodeResponseFlag)
	if err != nil {
		return nil, false, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, false, invalidFlag
	}

	iflagpercent, err := (*args.provider.store).GetNodeFlag(nodeA.nodePercentageFlag)
	if err != nil {
		return nil, false, err
	}

	flagpercent, ok := iflagpercent.(int32)
	if !ok {
		return nil, false, invalidPercentage
	}

	if flagpercent < rand.Int31n(int32(100)) {
		res, err := args.action()
		return res, true, err
	}

	switch flag {
	case scenario.Response_NORMAL:
		res, err := args.action()
		return res, true, err
	case scenario.Response_TIMEOUT:
		<-args.ctx.Done()
		return nil, true, nil
	case scenario.Response_ERROR:
		return nil, true, expectedError
	default:
		return nil, false, invalidResponse
	}
}

// podAndNodeResponse first calls podResponse and if pod response is not set calls nodeResponse
func podAndNodeResponse(args podNodeResponse) (interface{}, error) {
	pod, performedAction, err := podResponse(args.responseArgs, args.podResponseArgs)

	if err != nil {
		return nil, err
	}

	if performedAction {
		return pod, nil
	}

	node, performedAction, err := nodeResponse(args.responseArgs, args.nodeResponseArgs)

	if err != nil {
		return nil, err
	}

	if performedAction {
		return node, nil
	}

	return nil, nil
}
