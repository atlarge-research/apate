package cluster

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestDefault(t *testing.T) {
	clusterbuilder := Default()

	v := reflect.ValueOf(clusterbuilder)
	name := v.FieldByName("name")

	assert.Equal(t, name.String(), "Apate")
}

func TestWithName(t *testing.T) {
	clusterbuilder := New()

	clusterbuilder.WithName("Test")

	v := reflect.ValueOf(clusterbuilder)
	name := v.FieldByName("name")

	assert.Equal(t, name.String(), "Test")
}
