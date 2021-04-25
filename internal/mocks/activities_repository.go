// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repositories/activities_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	entities "github.com/edy4c7/darkpot-school-works/internal/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockActivitiesRepository is a mock of ActivitiesRepository interface
type MockActivitiesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockActivitiesRepositoryMockRecorder
}

// MockActivitiesRepositoryMockRecorder is the mock recorder for MockActivitiesRepository
type MockActivitiesRepositoryMockRecorder struct {
	mock *MockActivitiesRepository
}

// NewMockActivitiesRepository creates a new mock instance
func NewMockActivitiesRepository(ctrl *gomock.Controller) *MockActivitiesRepository {
	mock := &MockActivitiesRepository{ctrl: ctrl}
	mock.recorder = &MockActivitiesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockActivitiesRepository) EXPECT() *MockActivitiesRepositoryMockRecorder {
	return m.recorder
}

// GetAll mocks base method
func (m *MockActivitiesRepository) GetAll(arg0 context.Context) ([]*entities.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]*entities.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockActivitiesRepositoryMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockActivitiesRepository)(nil).GetAll), arg0)
}

// FindByUserID mocks base method
func (m *MockActivitiesRepository) FindByUserID(arg0 context.Context, arg1 string) ([]*entities.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", arg0, arg1)
	ret0, _ := ret[0].([]*entities.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID
func (mr *MockActivitiesRepositoryMockRecorder) FindByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockActivitiesRepository)(nil).FindByUserID), arg0, arg1)
}

// Create mocks base method
func (m *MockActivitiesRepository) Create(arg0 context.Context, arg1 *entities.Activity) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockActivitiesRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockActivitiesRepository)(nil).Create), arg0, arg1)
}
