package repository

import (
	"context"
	"errors"
	"testing"
	"time"

	pkgredis "github.com/AbdurrahmanA/short-url/pkg/redis"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

var (
	mockPkgRedis = pkgredis.Redis{}
	ctx          = context.TODO()
	mockRedis    redismock.ClientMock
	nilValue     = errors.New("Can not set nil value")
	ok           = "OK"
	notFound     = errors.New("Not Found")
	one          = int64(1)
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

	r := NewRedisRepository(&mockPkgRedis)
	for _, row := range mockRedisDatas {
		if row.err {
			mockRedis.ExpectGet(row.key).SetErr(redis.Nil)
		} else {
			mockRedis.ExpectGet(row.key).SetVal(row.value.(string))
		}

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
	close := setupRedis()
	defer close()

	r := NewRedisRepository(&mockPkgRedis)
	for _, row := range mockRedisDatas {
		if row.err {
			mockRedis.ExpectSet(row.key, row.value, 2*time.Second).SetErr(nilValue)
		} else {
			mockRedis.ExpectSet(row.key, row.value, 2*time.Second).SetVal(ok)
		}
		result, err := r.Set(row.key, row.value, 2*time.Second)
		if row.err {
			assert.Equal(t, nilValue, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, ok, result)
		}
	}
}

func TestDelete(t *testing.T) {
	close := setupRedis()
	defer close()

	r := NewRedisRepository(&mockPkgRedis)

	for _, row := range mockRedisDatas {
		if row.err {
			mockRedis.ExpectDel(row.key).SetErr(notFound)
		} else {
			mockRedis.ExpectDel(row.key).SetVal(one)
		}
		result, err := r.Delete(row.key)

		if row.err {
			assert.Equal(t, notFound, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, one, result)
		}
	}
}
