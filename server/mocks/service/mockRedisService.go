package service

import (
	"time"

	service "github.com/AbdurrahmanA/short-url/service"
)

type MockRedisService struct {
	setMethod    func(key string, value interface{}, ttl time.Duration) error
	getMethod    func(key string) (string, error)
	deleteMethod func(key string) error
}

func (m MockRedisService) Set(key string, value interface{}, ttl time.Duration) error {
	return m.setMethod(key, value, ttl)
}

func (m MockRedisService) Get(key string) (string, error) {
	return m.getMethod(key)
}

func (m MockRedisService) Delete(key string) error {
	return m.deleteMethod(key)
}

func NewMockRedisService(
	setMethod func(key string, value interface{}, ttl time.Duration) error,
	getMethod func(key string) (string, error),
	deleteMethod func(key string) error,
) service.IRedisService {
	mockRedisService := struct{ MockRedisService }{}
	mockRedisService.setMethod = setMethod
	mockRedisService.getMethod = getMethod
	mockRedisService.deleteMethod = deleteMethod
	return mockRedisService
}
