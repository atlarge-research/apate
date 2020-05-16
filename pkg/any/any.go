// Package any contains utilities for the any protobuf construct
// TODO remove this when moving node to CRD
package any

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// The type urls used to identify the marshalled types
const (
	int64URL     = "apate/int64"
	boolURL      = "apate/bool"
	responseURL  = "apate/response"
	podStatusURL = "apate/podstatus"
	stringURL    = "apate/string"
)

// MarshalOrDie marshals a value to Any or panics on failure
func MarshalOrDie(value interface{}) *anypb.Any {
	any, err := Marshal(value)
	if err == nil {
		return any
	}
	panic(err.Error())
}

// Marshal marshals a value to Any
func Marshal(value interface{}) (*anypb.Any, error) {
	// Use protoreflect.MessageDescriptor.FullName instead.
	var w proto.Message
	var t string
	switch v := value.(type) {
	case int:
		t = int64URL
		w = &wrapperspb.Int64Value{
			Value: int64(v),
		}
	case int32:
		t = int64URL
		w = &wrapperspb.Int64Value{
			Value: int64(v),
		}
	case int64:
		t = int64URL
		w = &wrapperspb.Int64Value{
			Value: v,
		}
	case scenario.Response:
		t = responseURL
		w = &wrapperspb.Int32Value{
			Value: int32(v),
		}
	case bool:
		t = boolURL
		w = &wrapperspb.BoolValue{
			Value: v,
		}
	case scenario.PodStatus:
		t = podStatusURL
		w = &wrapperspb.Int32Value{
			Value: int32(v),
		}
	case string:
		t = stringURL
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
		TypeUrl: t,
		Value:   val,
	}

	return any, nil
}

// UnmarshalOrDie calls Unmarshal but panics if an error occurs
func UnmarshalOrDie(any *anypb.Any) interface{} {
	val, err := Unmarshal(any)
	if err != nil {
		panic(err)
	}
	return val
}

// Unmarshal unmarshals an any struct into an interface
func Unmarshal(any *anypb.Any) (interface{}, error) {
	switch any.TypeUrl {
	case int64URL:
		res := &wrapperspb.Int64Value{}
		err := proto.Unmarshal(any.Value, res)
		if err != nil {
			return nil, err
		}
		return res.Value, nil
	case boolURL:
		res := wrapperspb.BoolValue{}
		err := proto.Unmarshal(any.Value, &res)
		if err != nil {
			return nil, err
		}
		return res.Value, nil
	case responseURL:
		res := wrapperspb.Int32Value{}
		err := proto.Unmarshal(any.Value, &res)
		if err != nil {
			return nil, err
		}
		return scenario.Response(res.Value), nil
	case podStatusURL:
		res := wrapperspb.Int32Value{}
		err := proto.Unmarshal(any.Value, &res)
		if err != nil {
			return nil, err
		}
		return scenario.PodStatus(res.Value), nil
	case stringURL:
		res := wrapperspb.StringValue{}
		err := proto.Unmarshal(any.Value, &res)
		if err != nil {
			return nil, err
		}
		return res.Value, nil
	default:
		return nil, fmt.Errorf("invalid type url: %s", any.TypeUrl)
	}
}
