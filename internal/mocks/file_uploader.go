// Code generated by MockGen. DO NOT EDIT.
// Source: internal/tools/file_uploader.go

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	multipart "mime/multipart"
	reflect "reflect"
)

// MockFileUploader is a mock of FileUploader interface
type MockFileUploader struct {
	ctrl     *gomock.Controller
	recorder *MockFileUploaderMockRecorder
}

// MockFileUploaderMockRecorder is the mock recorder for MockFileUploader
type MockFileUploaderMockRecorder struct {
	mock *MockFileUploader
}

// NewMockFileUploader creates a new mock instance
func NewMockFileUploader(ctrl *gomock.Controller) *MockFileUploader {
	mock := &MockFileUploader{ctrl: ctrl}
	mock.recorder = &MockFileUploaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFileUploader) EXPECT() *MockFileUploaderMockRecorder {
	return m.recorder
}

// Upload mocks base method
func (m *MockFileUploader) Upload(arg0 string, arg1 *multipart.FileHeader) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Upload indicates an expected call of Upload
func (mr *MockFileUploaderMockRecorder) Upload(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockFileUploader)(nil).Upload), arg0, arg1)
}

// Delete mocks base method
func (m *MockFileUploader) Delete(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockFileUploaderMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFileUploader)(nil).Delete), arg0)
}
