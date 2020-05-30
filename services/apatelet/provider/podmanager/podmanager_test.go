package podmanager

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestGetPodByName(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "a",
			Name:      "b",
			UID:       "c",
		},
	}

	pm.uidToPod["c"] = pod
	pm.nameToPod[getInternalPodName("a", "b")] = pod

	name, ok := pm.GetPodByName("a", "b")
	assert.Equal(t, pod, name)
	assert.True(t, ok)
}

func TestGetPodByNameNotFound(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	_, ok := pm.GetPodByName("a", "b")
	assert.False(t, ok)
}

func TestGetPodByUID(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "a",
			Name:      "b",
			UID:       "c",
		},
	}

	pm.uidToPod["c"] = pod
	pm.nameToPod[getInternalPodName("a", "b")] = pod

	name, ok := pm.GetPodByUID("c")
	assert.Equal(t, pod, name)
	assert.True(t, ok)
}

func TestGetPodByUIDNotFound(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	_, ok := pm.GetPodByUID("c")
	assert.False(t, ok)
}

func TestAddPod(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "a",
			Name:      "b",
			UID:       "c",
		},
	}

	pm.AddPod(pod)

	assert.EqualValues(t, pod, pm.nameToPod[getInternalPodName("a", "b")])
	assert.EqualValues(t, pod, pm.uidToPod["c"])
}

func TestDeletePod(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "a",
			Name:      "b",
			UID:       "c",
		},
	}

	pm.uidToPod["c"] = pod
	pm.nameToPod[getInternalPodName("a", "b")] = pod

	pm.DeletePod(pod)

	assert.Nil(t, pm.nameToPod[getInternalPodName("a", "b")])
	assert.Nil(t, pm.uidToPod["c"])
}

func TestDeletePodByName(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "a",
			Name:      "b",
			UID:       "c",
		},
	}

	pm.uidToPod["c"] = pod
	pm.nameToPod[getInternalPodName("a", "b")] = pod

	pm.DeletePodByName("a", "b")

	assert.Nil(t, pm.nameToPod[getInternalPodName("a", "b")])
	assert.Nil(t, pm.uidToPod["c"])
}

func TestGetAllPods(t *testing.T) {
	t.Parallel()

	pm := &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}

	pod1 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "a",
			Name:      "b",
			UID:       "c",
		},
	}

	pod2 := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "d",
			Name:      "e",
			UID:       "f",
		},
	}

	pm.nameToPod[getInternalPodName("a", "b")] = pod1
	pm.uidToPod["c"] = pod1

	pm.nameToPod[getInternalPodName("d", "e")] = pod2
	pm.uidToPod["f"] = pod2

	pods := pm.GetAllPods()

	assert.Contains(t, pods, pod1)
	assert.Contains(t, pods, pod2)
}
