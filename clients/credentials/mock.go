// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package credentials is a generated GoMock package.
package credentials

import (
	context "context"
	api "github.com/estafette/estafette-extension-gke/api"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Init mocks base method
func (m *MockClient) Init(ctx context.Context, paramsJSON, releaseName, credentialsPath string) (*api.GKECredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", ctx, paramsJSON, releaseName, credentialsPath)
	ret0, _ := ret[0].(*api.GKECredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Init indicates an expected call of Init
func (mr *MockClientMockRecorder) Init(ctx, paramsJSON, releaseName, credentialsPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockClient)(nil).Init), ctx, paramsJSON, releaseName, credentialsPath)
}

// GetCredentialsByName mocks base method
func (m *MockClient) GetCredentialsByName(c []api.GKECredentials, credentialName string) *api.GKECredentials {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentialsByName", c, credentialName)
	ret0, _ := ret[0].(*api.GKECredentials)
	return ret0
}

// GetCredentialsByName indicates an expected call of GetCredentialsByName
func (mr *MockClientMockRecorder) GetCredentialsByName(c, credentialName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentialsByName", reflect.TypeOf((*MockClient)(nil).GetCredentialsByName), c, credentialName)
}