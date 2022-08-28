package service

import (
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepo interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) *redis.StringCmd
	Delete(key string) error
}

type RedisService struct {
	redisShortURLRepo RedisRepo
}

func NewRedisService(redisShortURLRepo RedisRepo) *RedisService {
	return &RedisService{redisShortURLRepo: redisShortURLRepo}
}

func (r *RedisService) Set(key string, value interface{}, ttl time.Duration) error {
	err := r.redisShortURLRepo.Set(key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) Get(key string) (string, error) {
	cmd := r.redisShortURLRepo.Get(key)
	cmdErr := cmd.Err()
	if cmdErr != redis.Nil {
		return cmd.Val(), nil
	}

	return "", cmdErr
}

func (r *RedisService) Delete(key string) error {
	err := r.redisShortURLRepo.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
