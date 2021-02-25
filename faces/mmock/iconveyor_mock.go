// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/iostrovok/conveyor/faces (interfaces: IConveyor)

// Package mmock is a generated GoMock package.
package mmock

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	check "github.com/iostrovok/check"
	faces "github.com/iostrovok/conveyor/faces"
	nodes "github.com/iostrovok/conveyor/protobuf/go/nodes"
	reflect "reflect"
	time "time"
)

// MockIConveyor is a mock of IConveyor interface
type MockIConveyor struct {
	ctrl     *gomock.Controller
	recorder *MockIConveyorMockRecorder
}

// MockIConveyorMockRecorder is the mock recorder for MockIConveyor
type MockIConveyorMockRecorder struct {
	mock *MockIConveyor
}

// NewMockIConveyor creates a new mock instance
func NewMockIConveyor(ctrl *gomock.Controller) *MockIConveyor {
	mock := &MockIConveyor{ctrl: ctrl}
	mock.recorder = &MockIConveyorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIConveyor) EXPECT() *MockIConveyorMockRecorder {
	return m.recorder
}

// AddErrorHandler mocks base method
func (m *MockIConveyor) AddErrorHandler(arg0 faces.Name, arg1, arg2 int, arg3 faces.GiveBirth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddErrorHandler", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddErrorHandler indicates an expected call of AddErrorHandler
func (mr *MockIConveyorMockRecorder) AddErrorHandler(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddErrorHandler", reflect.TypeOf((*MockIConveyor)(nil).AddErrorHandler), arg0, arg1, arg2, arg3)
}

// AddFinalHandler mocks base method
func (m *MockIConveyor) AddFinalHandler(arg0 faces.Name, arg1, arg2 int, arg3 faces.GiveBirth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFinalHandler", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFinalHandler indicates an expected call of AddFinalHandler
func (mr *MockIConveyorMockRecorder) AddFinalHandler(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFinalHandler", reflect.TypeOf((*MockIConveyor)(nil).AddFinalHandler), arg0, arg1, arg2, arg3)
}

// AddHandler mocks base method
func (m *MockIConveyor) AddHandler(arg0 faces.Name, arg1, arg2 int, arg3 faces.GiveBirth) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddHandler", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddHandler indicates an expected call of AddHandler
func (mr *MockIConveyorMockRecorder) AddHandler(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHandler", reflect.TypeOf((*MockIConveyor)(nil).AddHandler), arg0, arg1, arg2, arg3)
}

// DefaultPriority mocks base method
func (m *MockIConveyor) DefaultPriority() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DefaultPriority")
	ret0, _ := ret[0].(int)
	return ret0
}

// DefaultPriority indicates an expected call of DefaultPriority
func (mr *MockIConveyorMockRecorder) DefaultPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DefaultPriority", reflect.TypeOf((*MockIConveyor)(nil).DefaultPriority))
}

// GetDefaultPriority mocks base method
func (m *MockIConveyor) GetDefaultPriority() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefaultPriority")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetDefaultPriority indicates an expected call of GetDefaultPriority
func (mr *MockIConveyorMockRecorder) GetDefaultPriority() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefaultPriority", reflect.TypeOf((*MockIConveyor)(nil).GetDefaultPriority))
}

// GetName mocks base method
func (m *MockIConveyor) GetName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetName indicates an expected call of GetName
func (mr *MockIConveyorMockRecorder) GetName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetName", reflect.TypeOf((*MockIConveyor)(nil).GetName))
}

// MetricPeriod mocks base method
func (m *MockIConveyor) MetricPeriod(arg0 time.Duration) faces.IConveyor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MetricPeriod", arg0)
	ret0, _ := ret[0].(faces.IConveyor)
	return ret0
}

// MetricPeriod indicates an expected call of MetricPeriod
func (mr *MockIConveyorMockRecorder) MetricPeriod(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MetricPeriod", reflect.TypeOf((*MockIConveyor)(nil).MetricPeriod), arg0)
}

// Run mocks base method
func (m *MockIConveyor) Run(arg0 faces.IInput) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run", arg0)
}

// Run indicates an expected call of Run
func (mr *MockIConveyorMockRecorder) Run(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockIConveyor)(nil).Run), arg0)
}

// RunRes mocks base method
func (m *MockIConveyor) RunRes(arg0 faces.IInput) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunRes", arg0)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunRes indicates an expected call of RunRes
func (mr *MockIConveyorMockRecorder) RunRes(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunRes", reflect.TypeOf((*MockIConveyor)(nil).RunRes), arg0)
}

// RunResTest mocks base method
func (m *MockIConveyor) RunResTest(arg0 faces.IInput, arg1 string) (interface{}, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunResTest", arg0, arg1)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RunResTest indicates an expected call of RunResTest
func (mr *MockIConveyorMockRecorder) RunResTest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunResTest", reflect.TypeOf((*MockIConveyor)(nil).RunResTest), arg0, arg1)
}

// RunTest mocks base method
func (m *MockIConveyor) RunTest(arg0 faces.IInput, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RunTest", arg0, arg1)
}

// RunTest indicates an expected call of RunTest
func (mr *MockIConveyorMockRecorder) RunTest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunTest", reflect.TypeOf((*MockIConveyor)(nil).RunTest), arg0, arg1)
}

// SetDefaultPriority mocks base method
func (m *MockIConveyor) SetDefaultPriority(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetDefaultPriority", arg0)
}

// SetDefaultPriority indicates an expected call of SetDefaultPriority
func (mr *MockIConveyorMockRecorder) SetDefaultPriority(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetDefaultPriority", reflect.TypeOf((*MockIConveyor)(nil).SetDefaultPriority), arg0)
}

// SetMasterNode mocks base method
func (m *MockIConveyor) SetMasterNode(arg0 string, arg1 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetMasterNode", arg0, arg1)
}

// SetMasterNode indicates an expected call of SetMasterNode
func (mr *MockIConveyorMockRecorder) SetMasterNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetMasterNode", reflect.TypeOf((*MockIConveyor)(nil).SetMasterNode), arg0, arg1)
}

// SetName mocks base method
func (m *MockIConveyor) SetName(arg0 string) faces.IConveyor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetName", arg0)
	ret0, _ := ret[0].(faces.IConveyor)
	return ret0
}

// SetName indicates an expected call of SetName
func (mr *MockIConveyorMockRecorder) SetName(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetName", reflect.TypeOf((*MockIConveyor)(nil).SetName), arg0)
}

// SetTestMode mocks base method
func (m *MockIConveyor) SetTestMode(arg0 bool, arg1 *check.C) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTestMode", arg0, arg1)
}

// SetTestMode indicates an expected call of SetTestMode
func (mr *MockIConveyorMockRecorder) SetTestMode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTestMode", reflect.TypeOf((*MockIConveyor)(nil).SetTestMode), arg0, arg1)
}

// SetTracer mocks base method
func (m *MockIConveyor) SetTracer(arg0 faces.ITrace, arg1 time.Duration) faces.IConveyor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetTracer", arg0, arg1)
	ret0, _ := ret[0].(faces.IConveyor)
	return ret0
}

// SetTracer indicates an expected call of SetTracer
func (mr *MockIConveyorMockRecorder) SetTracer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTracer", reflect.TypeOf((*MockIConveyor)(nil).SetTracer), arg0, arg1)
}

// SetWorkersCounter mocks base method
func (m *MockIConveyor) SetWorkersCounter(arg0 faces.IWorkersCounter) faces.IConveyor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetWorkersCounter", arg0)
	ret0, _ := ret[0].(faces.IConveyor)
	return ret0
}

// SetWorkersCounter indicates an expected call of SetWorkersCounter
func (mr *MockIConveyorMockRecorder) SetWorkersCounter(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetWorkersCounter", reflect.TypeOf((*MockIConveyor)(nil).SetWorkersCounter), arg0)
}

// Start mocks base method
func (m *MockIConveyor) Start(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockIConveyorMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockIConveyor)(nil).Start), arg0)
}

// Statistic mocks base method
func (m *MockIConveyor) Statistic() *nodes.SlaveNodeInfoRequest {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Statistic")
	ret0, _ := ret[0].(*nodes.SlaveNodeInfoRequest)
	return ret0
}

// Statistic indicates an expected call of Statistic
func (mr *MockIConveyorMockRecorder) Statistic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Statistic", reflect.TypeOf((*MockIConveyor)(nil).Statistic))
}

// Stop mocks base method
func (m *MockIConveyor) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop
func (mr *MockIConveyorMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockIConveyor)(nil).Stop))
}

// WaitAndStop mocks base method
func (m *MockIConveyor) WaitAndStop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WaitAndStop")
}

// WaitAndStop indicates an expected call of WaitAndStop
func (mr *MockIConveyorMockRecorder) WaitAndStop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WaitAndStop", reflect.TypeOf((*MockIConveyor)(nil).WaitAndStop))
}
