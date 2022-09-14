package service

import (
	"time"
)

type RedisRepo interface {
	Set(key string, value interface{}, ttl time.Duration) (string, error)
	Get(key string) (string, error)
	Delete(key string) (int64, error)
}

type RedisService struct {
	redisShortURLRepo RedisRepo
}

type IRedisService interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

func NewRedisService(redisShortURLRepo RedisRepo) IRedisService {
	return &RedisService{redisShortURLRepo: redisShortURLRepo}
}

func (r *RedisService) Set(key string, value interface{}, ttl time.Duration) error {
	_, err := r.redisShortURLRepo.Set(key, value, ttl)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) Get(key string) (string, error) {
	value, err := r.redisShortURLRepo.Get(key)
	if err != nil {
		return "", err
	}
	return value, err
}

func (r *RedisService) Delete(key string) error {
	_, err := r.redisShortURLRepo.Delete(key)
	if err != nil {
		return err
	}
	return nil
}
