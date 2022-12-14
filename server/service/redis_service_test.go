package service

import (
	"errors"
	"testing"
	"time"

	"github.com/AbdurrahmanA/short-url/mocks/repository"
	"github.com/go-redis/redis/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var mockRedisRepo *repository.MockRedisRepo
var mockRedisService IRedisService

type redisMockType struct {
	key   string
	value interface{}
	ttl   time.Duration
}

var RedisTestCases = []testCase{
	{arg1: "test", expected: "test", err: nil},
	{expected: "", err: redis.Nil},
}

func redisSetup(t *testing.T) func() {
	ct := gomock.NewController(t)
	defer ct.Finish()

	mockRedisRepo = repository.NewMockRedisRepo(ct)
	mockRedisService = NewRedisService(mockRedisRepo)
	return func() {
		mockRedisService = nil
		defer ct.Finish()
	}
}

func TestSetRedisWithoutError(t *testing.T) {
	td := redisSetup(t)
	defer td()

	test := redisMockType{key: "test", value: "test1", ttl: time.Minute}

	mockRedisRepo.EXPECT().Set(test.key, test.value, test.ttl).Return("OK", nil)
	err := mockRedisService.Set(test.key, test.value, test.ttl)
	assert.Equal(t, nil, err)
}

func TestSetRedisWithError(t *testing.T) {
	td := redisSetup(t)
	defer td()

	test := redisMockType{key: "test", value: "test1", ttl: time.Minute}

	mockRedisRepo.EXPECT().Set(test.key, test.value, test.ttl).DoAndReturn(func(key string, value interface{}, ttl time.Duration) (string, error) {
		return "", errors.New("something went wrong")
	})
	err := mockRedisService.Set(test.key, test.value, test.ttl)
	assert.Equal(t, errors.New("something went wrong"), err)
}

func TestGetRedis(t *testing.T) {
	td := redisSetup(t)
	defer td()

	for _, test := range RedisTestCases {
		mockRedisRepo.EXPECT().Get(test.arg1).Return(test.expected, test.err)
		result, err := mockRedisService.Get(test.arg1)

		if err != nil && test.err == nil {
			t.Error(err)
		}

		if test.err != nil {
			assert.Equal(t, test.err, err)
		}

		assert.Equal(t, test.expected, result)
	}
}

func TestDeleteRedis(t *testing.T) {
	td := redisSetup(t)
	defer td()

	for _, test := range RedisTestCases {
		mockRedisRepo.EXPECT().Delete(test.arg1).Return(int64(0), test.err)
		err := mockRedisService.Delete(test.arg1)

		assert.Equal(t, test.err, err)
	}
}
