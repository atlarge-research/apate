package condition

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestConditionReady(t *testing.T) {
	cond := New(true, corev1.NodeReady)

	// Make sure the ready messages are okay
	assert.Equal(t, readyMessage, cond.getMessage())
	assert.Equal(t, resourceReason, cond.getReason())

	// Set state to non-ready
	cond.Update(false)
	assert.Equal(t, notReadyMessage, cond.getMessage())
	assert.Equal(t, resourcePressureReason, cond.getReason())

	// Make sure the messages don't change if the state doesn't change
	cond.Update(false)
	assert.Equal(t, notReadyMessage, cond.getMessage())
	assert.Equal(t, resourcePressureReason, cond.getReason())
}

func TestConditionMemoryPressure(t *testing.T) {
	cond := New(true, corev1.NodeMemoryPressure)

	// Make sure the pressure messages are okay
	assert.Equal(t, notReadyMessage, cond.getMessage())
	assert.Equal(t, resourcePressureReason, cond.getReason())

	// Set state to non-pressure
	cond.Update(false)
	assert.Equal(t, readyMessage, cond.getMessage())
	assert.Equal(t, resourceReason, cond.getReason())

	// Make sure the messages don't change if the state doesn't change
	cond.Update(false)
	assert.Equal(t, readyMessage, cond.getMessage())
	assert.Equal(t, resourceReason, cond.getReason())
}
