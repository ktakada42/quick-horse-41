// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_login is a generated GoMock package.
package mock_login

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockrepositoryInterface is a mock of repositoryInterface interface.
type MockrepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockrepositoryInterfaceMockRecorder
}

// MockrepositoryInterfaceMockRecorder is the mock recorder for MockrepositoryInterface.
type MockrepositoryInterfaceMockRecorder struct {
	mock *MockrepositoryInterface
}

// NewMockrepositoryInterface creates a new mock instance.
func NewMockrepositoryInterface(ctrl *gomock.Controller) *MockrepositoryInterface {
	mock := &MockrepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockrepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockrepositoryInterface) EXPECT() *MockrepositoryInterfaceMockRecorder {
	return m.recorder
}

// CreateToken mocks base method.
func (m *MockrepositoryInterface) CreateToken(userId, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", userId, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockrepositoryInterfaceMockRecorder) CreateToken(userId, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockrepositoryInterface)(nil).CreateToken), userId, token)
}

// GetToken mocks base method.
func (m *MockrepositoryInterface) GetToken(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetToken indicates an expected call of GetToken.
func (mr *MockrepositoryInterfaceMockRecorder) GetToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockrepositoryInterface)(nil).GetToken), userId)
}

// UpdateToken mocks base method.
func (m *MockrepositoryInterface) UpdateToken(userId, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateToken", userId, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateToken indicates an expected call of UpdateToken.
func (mr *MockrepositoryInterfaceMockRecorder) UpdateToken(userId, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateToken", reflect.TypeOf((*MockrepositoryInterface)(nil).UpdateToken), userId, token)
}
