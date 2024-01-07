// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_review is a generated GoMock package.
package mock_review

import (
	review "app/review"
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

// GetReviews mocks base method.
func (m *MockrepositoryInterface) GetReviews(offset, limit int) ([]review.Review, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReviews", offset, limit)
	ret0, _ := ret[0].([]review.Review)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReviews indicates an expected call of GetReviews.
func (mr *MockrepositoryInterfaceMockRecorder) GetReviews(offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReviews", reflect.TypeOf((*MockrepositoryInterface)(nil).GetReviews), offset, limit)
}