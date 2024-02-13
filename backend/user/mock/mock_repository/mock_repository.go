// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// GetHashedPassword mocks base method.
func (m *MockUserRepository) GetHashedPassword(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHashedPassword", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHashedPassword indicates an expected call of GetHashedPassword.
func (mr *MockUserRepositoryMockRecorder) GetHashedPassword(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHashedPassword", reflect.TypeOf((*MockUserRepository)(nil).GetHashedPassword), userId)
}

// GetUserIdByEmail mocks base method.
func (m *MockUserRepository) GetUserIdByEmail(email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserIdByEmail", email)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserIdByEmail indicates an expected call of GetUserIdByEmail.
func (mr *MockUserRepositoryMockRecorder) GetUserIdByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserIdByEmail", reflect.TypeOf((*MockUserRepository)(nil).GetUserIdByEmail), email)
}