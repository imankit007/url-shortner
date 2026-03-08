package infrastructure

import (
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func NewRedisClient() *redis.Client {
	if RedisClient != nil {
		return RedisClient
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	return RedisClient
}

func InitializeRedisClient() {
	NewRedisClient()
}
