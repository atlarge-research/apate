// Package condition provides functions to more easily work with kubernetes node conditions
package condition

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	readyMessage    = "Apatelet is ready"
	notReadyMessage = "Apatelet is not ready"

	resourceReason         = "There are enough emulated resources"
	resourcePressureReason = "Emulated resources are (nearly) all used"
)

// Condition is used to maintain state over individual conditions
type Condition struct {
	val    bool
	negate bool

	condition  corev1.NodeConditionType
	transition metav1.Time
}

// New returns a new condition struct
func New(initial bool, condition corev1.NodeConditionType) Condition {
	return Condition{
		val:    initial,
		negate: condition != corev1.NodeReady,

		condition:  condition,
		transition: metav1.Now(),
	}
}

func (c *Condition) getStatus() corev1.ConditionStatus {
	if c.val {
		return corev1.ConditionTrue
	}

	return corev1.ConditionFalse
}

func (c *Condition) getMessage() string {
	if c.val != c.negate {
		return readyMessage
	}

	return notReadyMessage
}

func (c *Condition) getReason() string {
	if c.val != c.negate {
		return resourceReason
	}

	return resourcePressureReason
}

// Update returns the correct node condition based on the given condition
func (c *Condition) Update(condition bool) corev1.NodeCondition {
	// If the state is different, update state accordingly
	if c.val != condition {
		c.val = condition
		c.transition = metav1.Now()
	}

	return corev1.NodeCondition{
		Type:               c.condition,
		Status:             c.getStatus(),
		LastHeartbeatTime:  metav1.Now(),
		LastTransitionTime: c.transition,
		Reason:             c.getReason(),
		Message:            c.getMessage(),
	}
}

// Get returns the NodeCondition
func (c *Condition) Get() corev1.NodeCondition {
	return corev1.NodeCondition{
		Type:               c.condition,
		Status:             c.getStatus(),
		LastHeartbeatTime:  metav1.Now(),
		LastTransitionTime: c.transition,
		Reason:             c.getReason(),
		Message:            c.getMessage(),
	}
}
