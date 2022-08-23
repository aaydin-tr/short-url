package redis

import (
	"context"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type Redis struct {
	RDB *redis.Client
	CTX context.Context
}

var redisClient *redis.Client
var doOnce sync.Once

func NewRedisClient(URL string, Password string, Database int) *Redis {
	context := context.Background()

	doOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     URL,
			Password: Password,
			DB:       Database,
		})

		_, err := redisClient.Ping(context).Result()
		if err != nil {
			zap.S().Error("Error while connecting to Redis", err)
		}
	})
	zap.S().Infof("Redis connected successfully DB: %d", Database)
	return &Redis{
		RDB: redisClient,
		CTX: context,
	}
}

func (r *Redis) Ping() error {
	return r.RDB.Ping(r.CTX).Err()
}

func (r *Redis) Close() error {
	return r.RDB.Close()
}

func (r *Redis) Set(key string, value interface{}, ttl time.Duration) error {
	return r.RDB.Set(r.CTX, key, value, ttl).Err()
}

func (r *Redis) Get(key string) *redis.StringCmd {
	return r.RDB.Get(r.CTX, key)
}
