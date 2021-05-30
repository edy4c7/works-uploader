// Code generated by MockGen. DO NOT EDIT.
// Source: internal/services/activities_service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	entities "github.com/edy4c7/works-uploader/internal/entities"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockActivitiesService is a mock of ActivitiesService interface
type MockActivitiesService struct {
	ctrl     *gomock.Controller
	recorder *MockActivitiesServiceMockRecorder
}

// MockActivitiesServiceMockRecorder is the mock recorder for MockActivitiesService
type MockActivitiesServiceMockRecorder struct {
	mock *MockActivitiesService
}

// NewMockActivitiesService creates a new mock instance
func NewMockActivitiesService(ctrl *gomock.Controller) *MockActivitiesService {
	mock := &MockActivitiesService{ctrl: ctrl}
	mock.recorder = &MockActivitiesServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockActivitiesService) EXPECT() *MockActivitiesServiceMockRecorder {
	return m.recorder
}

// GetAll mocks base method
func (m *MockActivitiesService) GetAll(arg0 context.Context) ([]*entities.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", arg0)
	ret0, _ := ret[0].([]*entities.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockActivitiesServiceMockRecorder) GetAll(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockActivitiesService)(nil).GetAll), arg0)
}

// FindByUserID mocks base method
func (m *MockActivitiesService) FindByUserID(arg0 context.Context, arg1 string) ([]*entities.Activity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByUserID", arg0, arg1)
	ret0, _ := ret[0].([]*entities.Activity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByUserID indicates an expected call of FindByUserID
func (mr *MockActivitiesServiceMockRecorder) FindByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByUserID", reflect.TypeOf((*MockActivitiesService)(nil).FindByUserID), arg0, arg1)
}
