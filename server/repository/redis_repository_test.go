package repository

import (
	"context"
	"testing"

	pkgredis "github.com/AbdurrahmanA/short-url/pkg/redis"
	"github.com/go-redis/redis"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

var (
	mockPkgRedis = pkgredis.Redis{}
	ctx          = context.TODO()
	mockRedis    redismock.ClientMock
)

var mockRedisDatas = []struct {
	key   string
	value interface{}
	err   bool
}{
	{"a", "b", false},
	{"c", "d", false},
	{"e", nil, true},
}

func setupRedis() func() {
	db, mock := redismock.NewClientMock()
	mockPkgRedis = pkgredis.Redis{
		RDB: db,
		CTX: ctx,
	}
	mockRedis = mock
	return func() {
		mockRedis = nil
		db.Close()
	}
}

func TestGet(t *testing.T) {
	close := setupRedis()
	defer close()

	for _, row := range mockRedisDatas {
		if row.err {
			mockRedis.ExpectGet(row.key).SetErr(redis.Nil)
		} else {
			mockRedis.ExpectGet(row.key).SetVal(row.value.(string))
		}
	}

	r := NewRedisRepository(&mockPkgRedis)

	for _, row := range mockRedisDatas {
		val, err := r.Get(row.key)
		if row.err {
			assert.Equal(t, redis.Nil, err)
			assert.Empty(t, val)
		} else {
			assert.Equal(t, row.value, val)
			assert.NoError(t, err)
		}
	}
}

func TestSet(t *testing.T) {
}
