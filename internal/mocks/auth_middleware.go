// Code generated by MockGen. DO NOT EDIT.
// Source: internal/middlewares/auth_middleware.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockJWTMiddleware is a mock of JWTMiddleware interface
type MockJWTMiddleware struct {
	ctrl     *gomock.Controller
	recorder *MockJWTMiddlewareMockRecorder
}

// MockJWTMiddlewareMockRecorder is the mock recorder for MockJWTMiddleware
type MockJWTMiddlewareMockRecorder struct {
	mock *MockJWTMiddleware
}

// NewMockJWTMiddleware creates a new mock instance
func NewMockJWTMiddleware(ctrl *gomock.Controller) *MockJWTMiddleware {
	mock := &MockJWTMiddleware{ctrl: ctrl}
	mock.recorder = &MockJWTMiddlewareMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockJWTMiddleware) EXPECT() *MockJWTMiddlewareMockRecorder {
	return m.recorder
}

// CheckJWT mocks base method
func (m *MockJWTMiddleware) CheckJWT(w http.ResponseWriter, r *http.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckJWT", w, r)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckJWT indicates an expected call of CheckJWT
func (mr *MockJWTMiddlewareMockRecorder) CheckJWT(w, r interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckJWT", reflect.TypeOf((*MockJWTMiddleware)(nil).CheckJWT), w, r)
}
