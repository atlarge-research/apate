// Package podmanager contains the pod manager interface and a simple thread safe map-based implementation
package podmanager

import (
	"fmt"
	"sync"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// PodManager provides an opaque interface for thread safe fast pod mangement
type PodManager interface {
	// GetPodByName returns the pod identified by the namespace and name given in the parameters.
	GetPodByName(namespace string, name string) (*corev1.Pod, error)
	// GetPodByUID returns the pod identified by the uid given in the `uid` parameter.
	GetPodByUID(uid types.UID) (*corev1.Pod, error)
	// AddPod adds the pod specified in the `pod` parameter.
	AddPod(pod corev1.Pod)
	// DeletePod deletes the pod specified in the `pod` parameter.
	DeletePod(pod *corev1.Pod)
	// GetAllPods returns an array of all pods.
	GetAllPods() (ret []*corev1.Pod)
}

// podManager implements PodManager in a thread safe way using two maps and a RWLock
type podManager struct {
	uidToPod  map[types.UID]*corev1.Pod
	nameToPod map[string]*corev1.Pod
	lock      sync.RWMutex
}

// New creates a new PodManager fully initialised
func New() PodManager {
	return &podManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}
}

func (m *podManager) GetPodByName(namespace string, name string) (*corev1.Pod, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if pod, ok := m.nameToPod[getInternalPodName(namespace, name)]; ok {
		return pod, nil
	}
	return nil, fmt.Errorf("pod with namespace %v and name %v not found", namespace, name)
}

func (m *podManager) GetPodByUID(uid types.UID) (*corev1.Pod, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if pod, ok := m.uidToPod[uid]; ok {
		return pod, nil
	}

	return nil, fmt.Errorf("pod with uid %v not found", uid)
}

func (m *podManager) AddPod(pod corev1.Pod) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.uidToPod[pod.UID] = &pod
	m.nameToPod[getInternalPodName(pod.Namespace, pod.Name)] = &pod
}

func (m *podManager) DeletePod(pod *corev1.Pod) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.nameToPod, getInternalPodName(pod.Namespace, pod.Name))
	delete(m.uidToPod, pod.UID)
}

func (m *podManager) GetAllPods() (ret []*corev1.Pod) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, v := range m.uidToPod {
		ret = append(ret, v)
	}

	return
}

// getInternalPodName returns the concatenation of namespace and name which is used as an index inside the
// 	nameToPod map
func getInternalPodName(namespace string, name string) string {
	return namespace + name
}
