// Package any contains utilities for the any protobuf construct
// Don't look at this package. Thanks!
package any

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const int64Name = "int64"

// MarshalOrDie marshals a value to Any OR DIES!!! 666
func MarshalOrDie(value interface{}) *anypb.Any {
	any, err := Marshal(value)
	if err == nil {
		return any
	}
	panic(err.Error())
}

// Marshal marshals a value to Any
// TODO: improve @jonay2000 kthxbye
func Marshal(value interface{}) (*anypb.Any, error) {
	// Use protoreflect.MessageDescriptor.FullName instead.
	var w proto.Message
	var t string
	switch v := value.(type) {
	case int:
		t = int64Name
		w = &wrapperspb.Int64Value{
			Value: int64(v),
		}
	case int32:
		t = int64Name
		w = &wrapperspb.Int64Value{
			Value: int64(v),
		}
	case int64:
		t = int64Name
		w = &wrapperspb.Int64Value{
			Value: v,
		}
	case scenario.Response:
		t = "response"
		w = &wrapperspb.Int32Value{
			Value: int32(v),
		}
	case bool:
		t = "bool"
		w = &wrapperspb.BoolValue{
			Value: v,
		}
	case scenario.PodStatus:
		t = "podstatus"
		w = &wrapperspb.Int32Value{
			Value: int32(v),
		}
	case string:
		t = "string"
		w = &wrapperspb.StringValue{
			Value: v,
		}
	default:
		return nil, fmt.Errorf("unknown value type %v", v)
	}

	val, err := proto.Marshal(w)
	if err != nil {
		return nil, err
	}

	any := &anypb.Any{
		TypeUrl: "apate/" + t,
		Value:   val,
	}

	return any, nil
}
