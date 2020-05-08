package provider

import (
	"context"
	"math/rand"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/throw"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

const flagNotSetError = store.FlagNotSetError
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
func podResponse(args responseArgs, podA podResponseArgs) (interface{}, error) {
	iflag, err := (*args.provider.store).GetPodFlag(podA.name, podA.podResponseFlag)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, invalidFlag
	}

	iflagpercent, err := (*args.provider.store).GetPodFlag(podA.name, podA.podPercentageFlag)
	if err != nil {
		return nil, err
	}

	flagpercent, ok := iflagpercent.(int32)
	if !ok {
		return nil, invalidPercentage
	}

	if flagpercent == 0 {
		return nil, flagNotSetError
	}

	if flagpercent < rand.Int31n(int32(100)) {
		return args.action()
	}

	switch flag {
	case scenario.Response_NORMAL:
		return args.action()
	case scenario.Response_TIMEOUT:
		<-args.ctx.Done()
		return nil, nil
	case scenario.Response_ERROR:
		return nil, expectedError
	default:
		return nil, invalidResponse
	}
}

// nodeResponse checks the passed flags and calls the passed function on success
func nodeResponse(args responseArgs, nodeA nodeResponseArgs) (interface{}, error) {
	iflag, err := (*args.provider.store).GetNodeFlag(nodeA.nodeResponseFlag)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, invalidFlag
	}

	iflagpercent, err := (*args.provider.store).GetNodeFlag(nodeA.nodePercentageFlag)
	if err != nil {
		return nil, err
	}

	flagpercent, ok := iflagpercent.(int32)
	if !ok {
		return nil, invalidPercentage
	}

	if flagpercent < rand.Int31n(int32(100)) {
		return args.action()
	}

	switch flag {
	case scenario.Response_NORMAL:
		return args.action()
	case scenario.Response_TIMEOUT:
		<-args.ctx.Done()
		return nil, nil
	case scenario.Response_ERROR:
		return nil, expectedError
	default:
		return nil, invalidResponse
	}
}

// podAndNodeResponse first calls podResponse and if pod response is not set calls nodeResponse
func podAndNodeResponse(args podNodeResponse) (interface{}, error) {
	pod, err := podResponse(args.responseArgs, args.podResponseArgs)

	if err != flagNotSetError {
		return pod, err
	}

	return nodeResponse(args.responseArgs, args.nodeResponseArgs)
}
