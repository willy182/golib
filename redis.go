package golib

import (
	"crypto/tls"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

// redisClient variable for setting redis client
var redisClient = map[string]*redis.Client{}

// RedisClient function for setting redis client
func RedisClient(node string) *redis.Client {

	if val, ok := redisClient[node]; ok {
		return val
	}

	// get the real credentials from master connections
	host := os.Getenv(fmt.Sprintf("REDIS_%s_HOST", node))
	db, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("REDIS_%s_DB", node)))
	password := os.Getenv(fmt.Sprintf("REDIS_%s_PASS", node))
	maxRetries, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("REDIS_%s_MAX_RETRIES", node)))
	timeout, _ := strconv.Atoi(os.Getenv(fmt.Sprintf("REDIS_%s_IDLE_TIMEOUT", node)))
	idleTimeout := time.Duration(timeout)
	tlsSecured, _ := strconv.ParseBool(os.Getenv(fmt.Sprintf("REDIS_%s_TLS", node)))

	var conf *tls.Config

	if tlsSecured {
		conf = &tls.Config{
			InsecureSkipVerify: tlsSecured,
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:        host,
		Password:    password,
		DB:          db,
		MaxRetries:  maxRetries,
		IdleTimeout: time.Second * idleTimeout,
		TLSConfig:   conf,
	})

	redisClient[node] = client

	return client
}

// CloseRedis function for closing redis connection
func CloseRedis() {
	for _, c := range redisClient {
		c.Close()
	}
}
