// Code generated by MockGen. DO NOT EDIT.
// Source: ./service/redis_service.go

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockRedisRepo is a mock of RedisRepo interface.
type MockRedisRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRedisRepoMockRecorder
}

// MockRedisRepoMockRecorder is the mock recorder for MockRedisRepo.
type MockRedisRepoMockRecorder struct {
	mock *MockRedisRepo
}

// NewMockRedisRepo creates a new mock instance.
func NewMockRedisRepo(ctrl *gomock.Controller) *MockRedisRepo {
	mock := &MockRedisRepo{ctrl: ctrl}
	mock.recorder = &MockRedisRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRedisRepo) EXPECT() *MockRedisRepoMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockRedisRepo) Delete(key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", key)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRedisRepoMockRecorder) Delete(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRedisRepo)(nil).Delete), key)
}

// Get mocks base method.
func (m *MockRedisRepo) Get(key string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRedisRepoMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRedisRepo)(nil).Get), key)
}

// Set mocks base method.
func (m *MockRedisRepo) Set(key string, value interface{}, ttl time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", key, value, ttl)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockRedisRepoMockRecorder) Set(key, value, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockRedisRepo)(nil).Set), key, value, ttl)
}