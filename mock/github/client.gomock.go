// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/oinume/playground-go/mock/github (interfaces: Client)
//
// Generated by this command:
//
//	mockgen -destination=client.gomock.go -package=github . Client
//

// Package github is a generated GoMock package.
package github

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// ListBranches mocks base method.
func (m *MockClient) ListBranches(arg0 context.Context, arg1, arg2 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListBranches", arg0, arg1, arg2)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListBranches indicates an expected call of ListBranches.
func (mr *MockClientMockRecorder) ListBranches(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListBranches", reflect.TypeOf((*MockClient)(nil).ListBranches), arg0, arg1, arg2)
}
