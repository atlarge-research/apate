package provider

import (
	"context"
	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
	"github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	"math/rand"
)

const FlagNotSetError = store.FlagNotSetError

type magicArgs struct {
	ctx context.Context
	p *VKProvider
	action func() (interface{}, error)
}

type magicPodArgs struct {
	name string

	podResponseFlag events.PodEventFlag
	podPercentageFlag events.PodEventFlag
}

type magicNodeArgs struct {
	nodeResponseFlag events.NodeEventFlag
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
		return nil, wError("invalid flag type")
	}

	iflagp, err := (*args.p.store).GetPodFlag(podA.name, podA.podPercentageFlag)
	if err != nil {
		return nil, err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return nil, wError("invalid percentage type")
	}

	if flagp == 0 {
		return nil, FlagNotSetError
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
		return nil, wError("expected error")
	default:
		return nil, wError("invalid response")
	}
}

func magicNode(args magicArgs, nodeA magicNodeArgs) (interface{}, error) {
	iflag, err := (*args.p.store).GetNodeFlag(nodeA.nodeResponseFlag)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, wError("invalid flag type")
	}

	iflagp, err := (*args.p.store).GetNodeFlag(nodeA.nodePercentageFlag)
	if err != nil {
		return nil, err
	}

	flagp, ok := iflagp.(int32)
	if !ok {
		return nil, wError("invalid percentage type")
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
		return nil, wError("expected error")
	default:
		return nil, wError("invalid response")
	}
}

func magicPodAndNode(args magicPodNodeArgs) (interface{}, error) {
	pint, err := magicPod(args.magicArgs, args.magicPodArgs)

	if err != FlagNotSetError {
		return pint, err
	}

	return magicNode(args.magicArgs, args.magicNodeArgs)
}