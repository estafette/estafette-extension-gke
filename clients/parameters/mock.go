// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package parameters is a generated GoMock package.
package parameters

import (
	context "context"
	reflect "reflect"

	api "github.com/estafette/estafette-extension-gke/api"
	gomock "github.com/golang/mock/gomock"
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

// Init mocks base method.
func (m *MockClient) Init(ctx context.Context, paramsYAML string, credential *api.GKECredentials, gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName, releaseAction, releaseID string) (api.Params, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", ctx, paramsYAML, credential, gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName, releaseAction, releaseID)
	ret0, _ := ret[0].(api.Params)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Init indicates an expected call of Init.
func (mr *MockClientMockRecorder) Init(ctx, paramsYAML, credential, gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName, releaseAction, releaseID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockClient)(nil).Init), ctx, paramsYAML, credential, gitSource, gitOwner, gitName, appLabel, buildVersion, releaseName, releaseAction, releaseID)
}
