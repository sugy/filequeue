// Code generated by MockGen. DO NOT EDIT.
// Source: internal/command.go

// Package mock_filequeue is a generated GoMock package.
package mock_filequeue

import (
	exec "os/exec"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCommandExecutor is a mock of CommandExecutor interface.
type MockCommandExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockCommandExecutorMockRecorder
}

// MockCommandExecutorMockRecorder is the mock recorder for MockCommandExecutor.
type MockCommandExecutorMockRecorder struct {
	mock *MockCommandExecutor
}

// NewMockCommandExecutor creates a new mock instance.
func NewMockCommandExecutor(ctrl *gomock.Controller) *MockCommandExecutor {
	mock := &MockCommandExecutor{ctrl: ctrl}
	mock.recorder = &MockCommandExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommandExecutor) EXPECT() *MockCommandExecutorMockRecorder {
	return m.recorder
}

// Command mocks base method.
func (m *MockCommandExecutor) Command(name string, arg ...string) *exec.Cmd {
	m.ctrl.T.Helper()
	varargs := []interface{}{name}
	for _, a := range arg {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Command", varargs...)
	ret0, _ := ret[0].(*exec.Cmd)
	return ret0
}

// Command indicates an expected call of Command.
func (mr *MockCommandExecutorMockRecorder) Command(name interface{}, arg ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{name}, arg...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Command", reflect.TypeOf((*MockCommandExecutor)(nil).Command), varargs...)
}