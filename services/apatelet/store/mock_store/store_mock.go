// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store (interfaces: Store)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	store "github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store"
	gomock "github.com/golang/mock/gomock"
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

// GetNodeFlag mocks base method
func (m *MockStore) GetNodeFlag(arg0 int32) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeFlag", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeFlag indicates an expected call of GetNodeFlag
func (mr *MockStoreMockRecorder) GetNodeFlag(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeFlag", reflect.TypeOf((*MockStore)(nil).GetNodeFlag), arg0)
}

// GetPodFlag mocks base method
func (m *MockStore) GetPodFlag(arg0 string, arg1 int32) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPodFlag", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPodFlag indicates an expected call of GetPodFlag
func (mr *MockStoreMockRecorder) GetPodFlag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPodFlag", reflect.TypeOf((*MockStore)(nil).GetPodFlag), arg0, arg1)
}

// LenTasks mocks base method
func (m *MockStore) LenTasks() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LenTasks")
	ret0, _ := ret[0].(int)
	return ret0
}

// LenTasks indicates an expected call of LenTasks
func (mr *MockStoreMockRecorder) LenTasks() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LenTasks", reflect.TypeOf((*MockStore)(nil).LenTasks))
}

// PeekTask mocks base method
func (m *MockStore) PeekTask() (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PeekTask")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PeekTask indicates an expected call of PeekTask
func (mr *MockStoreMockRecorder) PeekTask() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PeekTask", reflect.TypeOf((*MockStore)(nil).PeekTask))
}

// PopTask mocks base method
func (m *MockStore) PopTask() (*store.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopTask")
	ret0, _ := ret[0].(*store.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PopTask indicates an expected call of PopTask
func (mr *MockStoreMockRecorder) PopTask() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopTask", reflect.TypeOf((*MockStore)(nil).PopTask))
}

// RemovePodTasks mocks base method
func (m *MockStore) RemovePodTasks(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePodTasks", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePodTasks indicates an expected call of RemovePodTasks
func (mr *MockStoreMockRecorder) RemovePodTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePodTasks", reflect.TypeOf((*MockStore)(nil).RemovePodTasks), arg0)
}

// SetNodeFlag mocks base method
func (m *MockStore) SetNodeFlag(arg0 int32, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNodeFlag", arg0, arg1)
}

// SetNodeFlag indicates an expected call of SetNodeFlag
func (mr *MockStoreMockRecorder) SetNodeFlag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNodeFlag", reflect.TypeOf((*MockStore)(nil).SetNodeFlag), arg0, arg1)
}

// SetNodeTasks mocks base method
func (m *MockStore) SetNodeTasks(arg0 []*store.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetNodeTasks", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetNodeTasks indicates an expected call of SetNodeTasks
func (mr *MockStoreMockRecorder) SetNodeTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNodeTasks", reflect.TypeOf((*MockStore)(nil).SetNodeTasks), arg0)
}

// SetPodFlag mocks base method
func (m *MockStore) SetPodFlag(arg0 string, arg1 int32, arg2 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetPodFlag", arg0, arg1, arg2)
}

// SetPodFlag indicates an expected call of SetPodFlag
func (mr *MockStoreMockRecorder) SetPodFlag(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPodFlag", reflect.TypeOf((*MockStore)(nil).SetPodFlag), arg0, arg1, arg2)
}

// SetPodTasks mocks base method
func (m *MockStore) SetPodTasks(arg0 string, arg1 []*store.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPodTasks", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPodTasks indicates an expected call of SetPodTasks
func (mr *MockStoreMockRecorder) SetPodTasks(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPodTasks", reflect.TypeOf((*MockStore)(nil).SetPodTasks), arg0, arg1)
}

// SetStartTime mocks base method
func (m *MockStore) SetStartTime(arg0 int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetStartTime", arg0)
}

// SetStartTime indicates an expected call of SetStartTime
func (mr *MockStoreMockRecorder) SetStartTime(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStartTime", reflect.TypeOf((*MockStore)(nil).SetStartTime), arg0)
}
