// Code generated by MockGen. DO NOT EDIT.
// Source: ./oneonone_dml_category.pb.go

// Package oneononeddlv1 is a generated GoMock package.
package oneononeddlv1

import (
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// MockCategoryRecorder is a mock of CategoryRecorder interface.
type MockCategoryRecorder struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryRecorderMockRecorder
}

// MockCategoryRecorderMockRecorder is the mock recorder for MockCategoryRecorder.
type MockCategoryRecorderMockRecorder struct {
	mock *MockCategoryRecorder
}

// NewMockCategoryRecorder creates a new mock instance.
func NewMockCategoryRecorder(ctrl *gomock.Controller) *MockCategoryRecorder {
	mock := &MockCategoryRecorder{ctrl: ctrl}
	mock.recorder = &MockCategoryRecorderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCategoryRecorder) EXPECT() *MockCategoryRecorderMockRecorder {
	return m.recorder
}

// FindByIDs mocks base method.
func (m *MockCategoryRecorder) FindByIDs(db *sql.DB, ids []uint64) ([]*Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDs", db, ids)
	ret0, _ := ret[0].([]*Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByIDs indicates an expected call of FindByIDs.
func (mr *MockCategoryRecorderMockRecorder) FindByIDs(db, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDs", reflect.TypeOf((*MockCategoryRecorder)(nil).FindByIDs), db, ids)
}

// Get mocks base method.
func (m *MockCategoryRecorder) Get(db *sql.DB, id uint64) (*Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", db, id)
	ret0, _ := ret[0].(*Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCategoryRecorderMockRecorder) Get(db, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCategoryRecorder)(nil).Get), db, id)
}

// List mocks base method.
func (m *MockCategoryRecorder) List(db *sql.DB, lastID *wrapperspb.UInt64Value, asc bool, limit int64) ([]*Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", db, lastID, asc, limit)
	ret0, _ := ret[0].([]*Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockCategoryRecorderMockRecorder) List(db, lastID, asc, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockCategoryRecorder)(nil).List), db, lastID, asc, limit)
}

// Save mocks base method.
func (m *MockCategoryRecorder) Save(db *sql.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", db)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockCategoryRecorderMockRecorder) Save(db interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockCategoryRecorder)(nil).Save), db)
}
