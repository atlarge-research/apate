package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	clusterbuilder := Default()

	assert.Equal(t, clusterbuilder.name, "Apate")
}

func TestWithName(t *testing.T) {
	clusterbuilder := New()

	clusterbuilder.WithName("Test")

	assert.Equal(t, clusterbuilder.name, "Test")
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
