// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repositories/users_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	entities "github.com/edy4c7/works-uploader/internal/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUsersRepository is a mock of UsersRepository interface
type MockUsersRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUsersRepositoryMockRecorder
}

// MockUsersRepositoryMockRecorder is the mock recorder for MockUsersRepository
type MockUsersRepositoryMockRecorder struct {
	mock *MockUsersRepository
}

// NewMockUsersRepository creates a new mock instance
func NewMockUsersRepository(ctrl *gomock.Controller) *MockUsersRepository {
	mock := &MockUsersRepository{ctrl: ctrl}
	mock.recorder = &MockUsersRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUsersRepository) EXPECT() *MockUsersRepositoryMockRecorder {
	return m.recorder
}

// Save mocks base method
func (m *MockUsersRepository) Save(arg0 context.Context, arg1 *entities.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockUsersRepositoryMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUsersRepository)(nil).Save), arg0, arg1)
}
