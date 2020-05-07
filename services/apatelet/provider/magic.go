package provider

import (
	"context"
	"math/rand"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
)

const flagNotSetError = store.FlagNotSetError
const expectedError = wError("expected error")
const invalidResponse = wError("invalid response")
const invalidPercentage = wError("invalid percentage type")
const invalidFlag = wError("invalid flag type")

type magicArgs struct {
	ctx    context.Context
	p      *VKProvider
	action func() (interface{}, error)
}

type magicPodArgs struct {
	name string

	podResponseFlag   events.PodEventFlag
	podPercentageFlag events.PodEventFlag
}

type magicNodeArgs struct {
	nodeResponseFlag   events.NodeEventFlag
	nodePercentageFlag events.NodeEventFlag
}

type magicPodNodeArgs struct {
	magicArgs
	magicPodArgs
	magicNodeArgs
}

func magicPod(args magicArgs, podA magicPodArgs) (interface{}, error) {
	iflag, err := (*args.p.store).GetPodFlag(podA.name, podA.podResponseFlag)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, invalidFlag
	}

	iflagp, err := (*args.p.store).GetPodFlag(podA.name, podA.podPercentageFlag)
	if err != nil {
		return nil, err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return nil, invalidPercentage
	}

	if flagp == 0 {
		return nil, flagNotSetError
	}

	if flagp < rand.Int31n(int32(100)) {
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

func magicNode(args magicArgs, nodeA magicNodeArgs) (interface{}, error) {
	iflag, err := (*args.p.store).GetNodeFlag(nodeA.nodeResponseFlag)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, invalidFlag
	}

	iflagp, err := (*args.p.store).GetNodeFlag(nodeA.nodePercentageFlag)
	if err != nil {
		return nil, err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return nil, invalidPercentage
	}

	if flagp < rand.Int31n(int32(100)) {
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

func magicPodAndNode(args magicPodNodeArgs) (interface{}, error) {
	pint, err := magicPod(args.magicArgs, args.magicPodArgs)

	if err != flagNotSetError {
		return pint, err
	}

	return magicNode(args.magicArgs, args.magicNodeArgs)
}
