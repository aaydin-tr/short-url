package repository

import (
	"context"
	"time"

	pkgredis "github.com/AbdurrahmanA/short-url/pkg/redis"
	"github.com/go-redis/redis/v8"
)

// go:generate mockgen -source=./repository/redis_repository.go -destination=./mocks/repository/mockRedisRepository.go -package=repository  RedisRepo
type RedisRepo interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) *redis.StringCmd
	Delete(key string) error
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

func (r *RedisRepository) Set(key string, value interface{}, ttl time.Duration) error {
	return r.rdb.Set(r.ctx, key, value, ttl).Err()
}

func (r *RedisRepository) Get(key string) *redis.StringCmd {
	return r.rdb.Get(r.ctx, key)
}

func (r *RedisRepository) Delete(key string) error {
	return r.rdb.Del(r.ctx, key).Err()
}
