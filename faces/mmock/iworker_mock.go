// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/iostrovok/conveyor/faces (interfaces: IWorker)

// Package mmock is a generated GoMock package.
package mmock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	faces "github.com/iostrovok/conveyor/faces"
	reflect "reflect"
)

// MockIWorker is a mock of IWorker interface
type MockIWorker struct {
	ctrl     *gomock.Controller
	recorder *MockIWorkerMockRecorder
}

// MockIWorkerMockRecorder is the mock recorder for MockIWorker
type MockIWorkerMockRecorder struct {
	mock *MockIWorker
}

// NewMockIWorker creates a new mock instance
func NewMockIWorker(ctrl *gomock.Controller) *MockIWorker {
	mock := &MockIWorker{ctrl: ctrl}
	mock.recorder = &MockIWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIWorker) EXPECT() *MockIWorkerMockRecorder {
	return m.recorder
}

// GetBorderCond mocks base method
func (m *MockIWorker) GetBorderCond() (faces.Name, faces.ManagerType, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBorderCond")
	ret0, _ := ret[0].(faces.Name)
	ret1, _ := ret[1].(faces.ManagerType)
	ret2, _ := ret[2].(bool)
	return ret0, ret1, ret2
}

// GetBorderCond indicates an expected call of GetBorderCond
func (mr *MockIWorkerMockRecorder) GetBorderCond() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBorderCond", reflect.TypeOf((*MockIWorker)(nil).GetBorderCond))
}

// ID mocks base method
func (m *MockIWorker) ID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (mr *MockIWorkerMockRecorder) ID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ID", reflect.TypeOf((*MockIWorker)(nil).ID))
}

// Name mocks base method
func (m *MockIWorker) Name() faces.Name {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(faces.Name)
	return ret0
}

// Name indicates an expected call of Name
func (mr *MockIWorkerMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockIWorker)(nil).Name))
}

// SetBorderCond mocks base method
func (m *MockIWorker) SetBorderCond(arg0 faces.ManagerType, arg1 bool, arg2 faces.Name) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetBorderCond", arg0, arg1, arg2)
}

// SetBorderCond indicates an expected call of SetBorderCond
func (mr *MockIWorkerMockRecorder) SetBorderCond(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBorderCond", reflect.TypeOf((*MockIWorker)(nil).SetBorderCond), arg0, arg1, arg2)
}

// Start mocks base method
func (m *MockIWorker) Start(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockIWorkerMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockIWorker)(nil).Start), arg0)
}

// Stop mocks base method
func (m *MockIWorker) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockIWorkerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockIWorker)(nil).Stop))
}
