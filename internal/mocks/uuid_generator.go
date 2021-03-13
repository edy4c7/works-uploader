// Code generated by MockGen. DO NOT EDIT.
// Source: internal/tools/uuid_generator.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUUIDGenerator is a mock of UUIDGenerator interface
type MockUUIDGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDGeneratorMockRecorder
}

// MockUUIDGeneratorMockRecorder is the mock recorder for MockUUIDGenerator
type MockUUIDGeneratorMockRecorder struct {
	mock *MockUUIDGenerator
}

// NewMockUUIDGenerator creates a new mock instance
func NewMockUUIDGenerator(ctrl *gomock.Controller) *MockUUIDGenerator {
	mock := &MockUUIDGenerator{ctrl: ctrl}
	mock.recorder = &MockUUIDGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUUIDGenerator) EXPECT() *MockUUIDGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockUUIDGenerator) Generate() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate")
	ret0, _ := ret[0].(string)
	return ret0
}

// Generate indicates an expected call of Generate
func (mr *MockUUIDGeneratorMockRecorder) Generate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockUUIDGenerator)(nil).Generate))
}
