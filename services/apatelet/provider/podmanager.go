package provider

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sync"
)

type PodManager struct {
	uidToPod      map[types.UID]*corev1.Pod
	nameToPod	  map[string]*corev1.Pod
	lock  	  sync.RWMutex
}

func getInternalPodName(namespace string, name string) string {
	return namespace + name
}

func NewPodManager() PodManager {
	return PodManager{
		uidToPod:  make(map[types.UID]*corev1.Pod),
		nameToPod: make(map[string]*corev1.Pod),
	}
}

func (m *PodManager) GetPodByName(namespace string, name string) *corev1.Pod{
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.nameToPod[getInternalPodName(namespace, name)]
}

func (m *PodManager) GetPodByUID(uid types.UID) *corev1.Pod {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.uidToPod[uid]
}

func (m *PodManager) AddPod(pod corev1.Pod) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.uidToPod[pod.UID] = &pod
	m.nameToPod[getInternalPodName(pod.Namespace, pod.Name)] = &pod
}

func (m *PodManager) DeletePod(pod *corev1.Pod) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.nameToPod, getInternalPodName(pod.Namespace, pod.Name))
	delete(m.uidToPod, pod.UID)
}

func (m *PodManager) GetAllPods() (ret []*corev1.Pod) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	for _, v := range m.uidToPod {
		ret = append(ret, v)
	}

	return
}
