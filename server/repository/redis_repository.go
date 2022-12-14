package repository

import (
	"context"
	"time"

	pkgredis "github.com/AbdurrahmanA/short-url/pkg/redis"
	"github.com/go-redis/redis/v8"
)

// go:generate mockgen -source=./repository/redis_repository.go -destination=./mocks/repository/mockRedisRepository.go -package=repository  RedisRepo
type RedisRepo interface {
	Set(key string, value interface{}, ttl time.Duration) (string, error)
	Get(key string) (string, error)
	Delete(key string) (int64, error)
}

type RedisRepository struct {
	rdb *redis.Client
	ctx context.Context
}

func NewRedisRepository(redis *pkgredis.Redis) *RedisRepository {
	return &RedisRepository{
		rdb: redis.RDB,
		ctx: redis.CTX,
	}
}

func (r *RedisRepository) Set(key string, value interface{}, ttl time.Duration) (string, error) {
	return r.rdb.Set(r.ctx, key, value, ttl).Result()
}

func (r *RedisRepository) Get(key string) (string, error) {
	return r.rdb.Get(r.ctx, key).Result()
}

func (r *RedisRepository) Delete(key string) (int64, error) {
	return r.rdb.Del(r.ctx, key).Result()
}
