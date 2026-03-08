package infrastructure

import (
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}
