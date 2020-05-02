// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/atlarge-research/opendc-emulate-kubernetes/services/apatelet/store (interfaces: Store)

// Package mock_store is a generated GoMock package.
package mock_store

import (
	apatelet "github.com/atlarge-research/opendc-emulate-kubernetes/api/apatelet"
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

// EnqueueTasks mocks base method
func (m *MockStore) EnqueueTasks(arg0 []*apatelet.Task) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "EnqueueTasks", arg0)
}

// EnqueueTasks indicates an expected call of EnqueueTasks
func (mr *MockStoreMockRecorder) EnqueueTasks(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnqueueTasks", reflect.TypeOf((*MockStore)(nil).EnqueueTasks), arg0)
}

// GetFlag mocks base method
func (m *MockStore) GetFlag(arg0 int32) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlag", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFlag indicates an expected call of GetFlag
func (mr *MockStoreMockRecorder) GetFlag(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlag", reflect.TypeOf((*MockStore)(nil).GetFlag), arg0)
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
func (m *MockStore) PopTask() (*apatelet.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PopTask")
	ret0, _ := ret[0].(*apatelet.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PopTask indicates an expected call of PopTask
func (mr *MockStoreMockRecorder) PopTask() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PopTask", reflect.TypeOf((*MockStore)(nil).PopTask))
}

// SetFlag mocks base method
func (m *MockStore) SetFlag(arg0 int32, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetFlag", arg0, arg1)
}

// SetFlag indicates an expected call of SetFlag
func (mr *MockStoreMockRecorder) SetFlag(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetFlag", reflect.TypeOf((*MockStore)(nil).SetFlag), arg0, arg1)
}
