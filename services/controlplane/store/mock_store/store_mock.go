// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store (interfaces: Store)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	apatelet "github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
	health "github.com/atlarge-research/opendc-emulate-kubernetes/api/health"
	kubeconfig "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/kubernetes/kubeconfig"
	scenario "github.com/atlarge-research/opendc-emulate-kubernetes/pkg/scenario"
	store "github.com/atlarge-research/opendc-emulate-kubernetes/services/controlplane/store"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	reflect "reflect"
)

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddNode mocks base method
func (m *MockStore) AddNode(arg0 *store.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNode indicates an expected call of AddNode
func (mr *MockStoreMockRecorder) AddNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNode", reflect.TypeOf((*MockStore)(nil).AddNode), arg0)
}

// AddResourcesToQueue mocks base method
func (m *MockStore) AddResourcesToQueue(arg0 []scenario.NodeResources) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddResourcesToQueue", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddResourcesToQueue indicates an expected call of AddResourcesToQueue
func (mr *MockStoreMockRecorder) AddResourcesToQueue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddResourcesToQueue", reflect.TypeOf((*MockStore)(nil).AddResourcesToQueue), arg0)
}

// ClearNodes mocks base method
func (m *MockStore) ClearNodes() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearNodes")
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearNodes indicates an expected call of ClearNodes
func (mr *MockStoreMockRecorder) ClearNodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearNodes", reflect.TypeOf((*MockStore)(nil).ClearNodes))
}

// GetApateletScenario mocks base method
func (m *MockStore) GetApateletScenario() (*apatelet.ApateletScenario, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApateletScenario")
	ret0, _ := ret[0].(*apatelet.ApateletScenario)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApateletScenario indicates an expected call of GetApateletScenario
func (mr *MockStoreMockRecorder) GetApateletScenario() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApateletScenario", reflect.TypeOf((*MockStore)(nil).GetApateletScenario))
}

// GetKubeConfig mocks base method
func (m *MockStore) GetKubeConfig() (kubeconfig.KubeConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKubeConfig")
	ret0, _ := ret[0].(kubeconfig.KubeConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKubeConfig indicates an expected call of GetKubeConfig
func (mr *MockStoreMockRecorder) GetKubeConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKubeConfig", reflect.TypeOf((*MockStore)(nil).GetKubeConfig))
}

// GetNode mocks base method
func (m *MockStore) GetNode(arg0 uuid.UUID) (store.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNode", arg0)
	ret0, _ := ret[0].(store.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNode indicates an expected call of GetNode
func (mr *MockStoreMockRecorder) GetNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNode", reflect.TypeOf((*MockStore)(nil).GetNode), arg0)
}

// GetNodes mocks base method
func (m *MockStore) GetNodes() ([]store.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes")
	ret0, _ := ret[0].([]store.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes
func (mr *MockStoreMockRecorder) GetNodes() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockStore)(nil).GetNodes))
}

// GetNodesBySelector mocks base method
func (m *MockStore) GetNodesBySelector(arg0 string) ([]store.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodesBySelector", arg0)
	ret0, _ := ret[0].([]store.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodesBySelector indicates an expected call of GetNodesBySelector
func (mr *MockStoreMockRecorder) GetNodesBySelector(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodesBySelector", reflect.TypeOf((*MockStore)(nil).GetNodesBySelector), arg0)
}

// GetResourceFromQueue mocks base method
func (m *MockStore) GetResourceFromQueue() (*scenario.NodeResources, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResourceFromQueue")
	ret0, _ := ret[0].(*scenario.NodeResources)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetResourceFromQueue indicates an expected call of GetResourceFromQueue
func (mr *MockStoreMockRecorder) GetResourceFromQueue() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResourceFromQueue", reflect.TypeOf((*MockStore)(nil).GetResourceFromQueue))
}

// RemoveNode mocks base method
func (m *MockStore) RemoveNode(arg0 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNode indicates an expected call of RemoveNode
func (mr *MockStoreMockRecorder) RemoveNode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNode", reflect.TypeOf((*MockStore)(nil).RemoveNode), arg0)
}

// RemoveNodes mocks base method
func (m *MockStore) RemoveNodes(arg0 []uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveNodes", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveNodes indicates an expected call of RemoveNodes
func (mr *MockStoreMockRecorder) RemoveNodes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNodes", reflect.TypeOf((*MockStore)(nil).RemoveNodes), arg0)
}

// SetApateletScenario mocks base method
func (m *MockStore) SetApateletScenario(arg0 *apatelet.ApateletScenario) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetApateletScenario", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetApateletScenario indicates an expected call of SetApateletScenario
func (mr *MockStoreMockRecorder) SetApateletScenario(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetApateletScenario", reflect.TypeOf((*MockStore)(nil).SetApateletScenario), arg0)
}

// SetKubeConfig mocks base method
func (m *MockStore) SetKubeConfig(arg0 kubeconfig.KubeConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetKubeConfig", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetKubeConfig indicates an expected call of SetKubeConfig
func (mr *MockStoreMockRecorder) SetKubeConfig(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetKubeConfig", reflect.TypeOf((*MockStore)(nil).SetKubeConfig), arg0)
}

// SetNodeStatus mocks base method
func (m *MockStore) SetNodeStatus(arg0 uuid.UUID, arg1 health.Status) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNodeStatus", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNodeStatus indicates an expected call of SetNodeStatus
func (mr *MockStoreMockRecorder) SetNodeStatus(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNodeStatus", reflect.TypeOf((*MockStore)(nil).SetNodeStatus), arg0, arg1)
}
