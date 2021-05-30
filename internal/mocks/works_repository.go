// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repositories/works_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	entities "github.com/edy4c7/works-uploader/internal/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockWorksRepository is a mock of WorksRepository interface
type MockWorksRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWorksRepositoryMockRecorder
}

// MockWorksRepositoryMockRecorder is the mock recorder for MockWorksRepository
type MockWorksRepositoryMockRecorder struct {
	mock *MockWorksRepository
}

// NewMockWorksRepository creates a new mock instance
func NewMockWorksRepository(ctrl *gomock.Controller) *MockWorksRepository {
	mock := &MockWorksRepository{ctrl: ctrl}
	mock.recorder = &MockWorksRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockWorksRepository) EXPECT() *MockWorksRepositoryMockRecorder {
	return m.recorder
}

// GetAll mocks base method
func (m *MockWorksRepository) GetAll(arg0 context.Context) ([]*entities.Work, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]*entities.Work)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockWorksRepositoryMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockWorksRepository)(nil).GetAll), arg0)
}

// FindByID mocks base method
func (m *MockWorksRepository) FindByID(arg0 context.Context, arg1 uint64) (*entities.Work, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*entities.Work)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID
func (mr *MockWorksRepositoryMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockWorksRepository)(nil).FindByID), arg0, arg1)
}

// Save mocks base method
func (m *MockWorksRepository) Save(arg0 context.Context, arg1 *entities.Work) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockWorksRepositoryMockRecorder) Save(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockWorksRepository)(nil).Save), arg0, arg1)
}

// DeleteByID mocks base method
func (m *MockWorksRepository) DeleteByID(arg0 context.Context, arg1 uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID
func (mr *MockWorksRepositoryMockRecorder) DeleteByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockWorksRepository)(nil).DeleteByID), arg0, arg1)
}
