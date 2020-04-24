package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	clusterbuilder := Default()

	assert.Equal(t, "Apate", clusterbuilder.name)
}

func TestWithName(t *testing.T) {
	clusterbuilder := New()

	clusterbuilder.WithName("Test")

	assert.Equal(t, "Test", clusterbuilder.name)
}

func TestEmptyName(t *testing.T) {
	clusterbuilder := New()

	clusterbuilder.WithName("")

	_, err := clusterbuilder.Create()
	assert.Error(t, err)
}

func TestEmptyNameForce(t *testing.T) {
	clusterbuilder := New()

	clusterbuilder.WithName("")

	_, err := clusterbuilder.ForceCreate()
	assert.Error(t, err)
}
