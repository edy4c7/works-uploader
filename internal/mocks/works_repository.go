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
func (m *MockWorksRepository) GetAll(ctx context.Context, offset, limit int) ([]*entities.Work, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, offset, limit)
	ret0, _ := ret[0].([]*entities.Work)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockWorksRepositoryMockRecorder) GetAll(ctx, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockWorksRepository)(nil).GetAll), ctx, offset, limit)
}

// CountAll mocks base method
func (m *MockWorksRepository) CountAll(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountAll", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountAll indicates an expected call of CountAll
func (mr *MockWorksRepositoryMockRecorder) CountAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountAll", reflect.TypeOf((*MockWorksRepository)(nil).CountAll), arg0)
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

// Create mocks base method
func (m *MockWorksRepository) Create(arg0 context.Context, arg1 *entities.Work) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockWorksRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockWorksRepository)(nil).Create), arg0, arg1)
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
