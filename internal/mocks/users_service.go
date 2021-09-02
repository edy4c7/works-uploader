// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/users_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	beans "github.com/edy4c7/works-uploader/internal/beans"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUsersService is a mock of UsersService interface
type MockUsersService struct {
	ctrl     *gomock.Controller
	recorder *MockUsersServiceMockRecorder
}

// MockUsersServiceMockRecorder is the mock recorder for MockUsersService
type MockUsersServiceMockRecorder struct {
	mock *MockUsersService
}

// NewMockUsersService creates a new mock instance
func NewMockUsersService(ctrl *gomock.Controller) *MockUsersService {
	mock := &MockUsersService{ctrl: ctrl}
	mock.recorder = &MockUsersServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsersService) EXPECT() *MockUsersServiceMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockUsersService) Save(arg0 context.Context, arg1 *beans.UserFormBean) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockUsersServiceMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUsersService)(nil).Save), arg0, arg1)
}