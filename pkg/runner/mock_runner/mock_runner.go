// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/atlarge-research/apate/pkg/runner (interfaces: ApateletRunner)

// Package mock_runner is a generated GoMock package.
package mock_runner

import (
	context "context"
	env "github.com/atlarge-research/apate/pkg/env"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockApateletRunner is a mock of ApateletRunner interface
type MockApateletRunner struct {
	ctrl     *gomock.Controller
	recorder *MockApateletRunnerMockRecorder
}

// MockApateletRunnerMockRecorder is the mock recorder for MockApateletRunner
type MockApateletRunnerMockRecorder struct {
	mock *MockApateletRunner
}

// NewMockApateletRunner creates a new mock instance
func NewMockApateletRunner(ctrl *gomock.Controller) *MockApateletRunner {
	mock := &MockApateletRunner{ctrl: ctrl}
	mock.recorder = &MockApateletRunnerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockApateletRunner) EXPECT() *MockApateletRunnerMockRecorder {
	return m.recorder
}

// SpawnApatelets mocks base method
func (m *MockApateletRunner) SpawnApatelets(arg0 context.Context, arg1 int, arg2 env.ApateletEnvironment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SpawnApatelets", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SpawnApatelets indicates an expected call of SpawnApatelets
func (mr *MockApateletRunnerMockRecorder) SpawnApatelets(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SpawnApatelets", reflect.TypeOf((*MockApateletRunner)(nil).SpawnApatelets), arg0, arg1, arg2)
}
