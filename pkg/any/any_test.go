package any

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/atlarge-research/opendc-emulate-kubernetes/api/scenario"
)

func TestMarshalUnmarshalInt(t *testing.T) {
	any, err := Marshal(42)
	assert.NoError(t, err)
	assert.Equal(t, int64URL, any.TypeUrl)

	val, err := Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, int64(42), val)
}

func TestMarshalUnmarshalInt64(t *testing.T) {
	any, err := Marshal(int64(42))
	assert.NoError(t, err)
	assert.Equal(t, int64URL, any.TypeUrl)

	val, err := Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, int64(42), val)
}

func TestMarshalUnmarshalInt32(t *testing.T) {
	any, err := Marshal(int32(42))
	assert.NoError(t, err)
	assert.Equal(t, int64URL, any.TypeUrl)

	val, err := Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, int64(42), val)
}

func TestMarshalUnmarshalBool(t *testing.T) {
	any, err := Marshal(true)
	assert.NoError(t, err)
	assert.Equal(t, boolURL, any.TypeUrl)

	val, err := Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, true, val)

	any, err = Marshal(false)
	assert.NoError(t, err)
	assert.Equal(t, boolURL, any.TypeUrl)

	val, err = Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, false, val)
}

func TestMarshalUnmarshalResponse(t *testing.T) {
	any, err := Marshal(scenario.Response_RESPONSE_ERROR)
	assert.NoError(t, err)
	assert.Equal(t, responseURL, any.TypeUrl)

	val, err := Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, scenario.Response_RESPONSE_ERROR, val)
}

func TestMarshalUnmarshalPodStatus(t *testing.T) {
	any, err := Marshal(scenario.PodStatus_POD_STATUS_SUCCEEDED)
	assert.NoError(t, err)
	assert.Equal(t, podStatusURL, any.TypeUrl)

	val, err := Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, scenario.PodStatus_POD_STATUS_SUCCEEDED, val)
}

func TestMarshalUnmarshalString(t *testing.T) {
	str := "testing string"

	any, err := Marshal(str)
	assert.NoError(t, err)
	assert.Equal(t, stringURL, any.TypeUrl)

	val, err := Unmarshal(any)
	assert.NoError(t, err)

	assert.Equal(t, str, val)
}

func TestUnmarshalInvalidUrl(t *testing.T) {
	any := anypb.Any{
		TypeUrl: "invalid-url",
		Value:   nil,
	}

	v, err := Unmarshal(&any)
	assert.Error(t, err)
	assert.Nil(t, v)
}

func TestMarshalInvalidType(t *testing.T) {
	type invalidType bool
	v, err := Marshal(invalidType(true))
	assert.Error(t, err)
	assert.Nil(t, v)
}
