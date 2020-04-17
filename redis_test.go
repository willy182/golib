package golib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/redis.v4"
)

func TestRedisClient(t *testing.T) {
	client := redis.NewClient(&redis.Options{})
	redisClient = make(map[string]*redis.Client)
	redisClient["test"] = client

	t.Run("OK NODE RedisClient", func(t *testing.T) {
		assert.Equal(t, client, RedisClient("test"))
	})

	t.Run("NOK NODE RedisClient", func(t *testing.T) {
		assert.NotEqual(t, client, RedisClient("new"))
	})
}

func TestCloseRedis(t *testing.T) {
	client := redis.NewClient(&redis.Options{})
	redisClient = make(map[string]*redis.Client)
	redisClient["test"] = client

	t.Run("OK CloseRedis", func(t *testing.T) {
		CloseRedis()
	})
}
