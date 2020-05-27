package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	t.Parallel()

	clusterbuilder := Default()

	assert.Equal(t, "Apate", clusterbuilder.name)
}

func TestWithName(t *testing.T) {
	t.Parallel()

	clusterbuilder := New()

	clusterbuilder.WithName("Test")

	assert.Equal(t, "Test", clusterbuilder.name)
}

func TestEmptyName(t *testing.T) {
	t.Parallel()

	clusterbuilder := New()

	clusterbuilder.WithName("")

	_, err := clusterbuilder.Create()
	assert.Error(t, err)
}

func TestEmptyNameForce(t *testing.T) {
	t.Parallel()

	clusterbuilder := New()

	clusterbuilder.WithName("")

	_, err := clusterbuilder.ForceCreate()
	assert.Error(t, err)
}
