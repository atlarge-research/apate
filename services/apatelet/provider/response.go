package provider

import (
	"context"
	"errors"
	"math/rand"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
	"github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario/events"
)

type responseArgs struct {
	ctx      context.Context
	provider *Provider

	nodeResponseFlag   events.NodeEventFlag
	nodePercentageFlag events.NodeEventFlag

	action   func() (interface{}, error)
}

// nodeResponse checks the passed flags and calls the passed function on success
func nodeResponse(args responseArgs) (interface{}, error) {
	iflag, err := (*args.provider.store).GetNodeFlag(args.nodeResponseFlag)
	if err != nil {
		return nil, err
	}

	flag, ok := iflag.(scenario.Response)
	if !ok {
		return nil, errors.New("nodeResponse couldn't cast flag to response")
	}

	iflagpercent, err := (*args.provider.store).GetNodeFlag(args.nodePercentageFlag)
	if err != nil {
		return nil, err
	}

	flagpercent, ok := iflagpercent.(int32)
	if !ok {
		return nil, errors.New("nodeResponse couldn't cast percent to int")
	}

	if flagpercent < rand.Int31n(int32(100)) {
		res, err := args.action()
		return res, err
	}

	switch flag {
	case scenario.Response_NORMAL:
		res, err := args.action()
		return res, err
	case scenario.Response_TIMEOUT:
		<-args.ctx.Done()
		return nil, nil
	case scenario.Response_ERROR:
		return nil, errors.New("nodeResponse expected error")
	default:
		return nil, errors.New("nodeResponse invalid scenario")
	}
}
