// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GGjahoon/IZH/application/user/rpc/internal/model (interfaces: UserModel)

// Package mockuserModel is a generated GoMock package.
package mockuserModel

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	model "github.com/GGjahoon/IZH/application/user/rpc/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockUserModel is a mock of UserModel interface.
type MockUserModel struct {
	ctrl     *gomock.Controller
	recorder *MockUserModelMockRecorder
}

// MockUserModelMockRecorder is the mock recorder for MockUserModel.
type MockUserModelMockRecorder struct {
	mock *MockUserModel
}

// NewMockUserModel creates a new mock instance.
func NewMockUserModel(ctrl *gomock.Controller) *MockUserModel {
	mock := &MockUserModel{ctrl: ctrl}
	mock.recorder = &MockUserModelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserModel) EXPECT() *MockUserModelMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUserModel) Delete(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUserModelMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserModel)(nil).Delete), arg0, arg1)
}

// FindByMobile mocks base method.
func (m *MockUserModel) FindByMobile(arg0 context.Context, arg1 string) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByMobile", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByMobile indicates an expected call of FindByMobile.
func (mr *MockUserModelMockRecorder) FindByMobile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByMobile", reflect.TypeOf((*MockUserModel)(nil).FindByMobile), arg0, arg1)
}

// FindOne mocks base method.
func (m *MockUserModel) FindOne(arg0 context.Context, arg1 int64) (*model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOne", arg0, arg1)
	ret0, _ := ret[0].(*model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOne indicates an expected call of FindOne.
func (mr *MockUserModelMockRecorder) FindOne(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOne", reflect.TypeOf((*MockUserModel)(nil).FindOne), arg0, arg1)
}

// Insert mocks base method.
func (m *MockUserModel) Insert(arg0 context.Context, arg1 *model.User) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockUserModelMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockUserModel)(nil).Insert), arg0, arg1)
}

// Update mocks base method.
func (m *MockUserModel) Update(arg0 context.Context, arg1 *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockUserModelMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUserModel)(nil).Update), arg0, arg1)
}
