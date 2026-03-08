package counter

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisCounterKey = "url-shortener:global-counter"

type URLCodeCounter struct {
	redisClient *redis.Client
}

func NewURLCodeCounter(redisClient *redis.Client) *URLCodeCounter {
	return &URLCodeCounter{
		redisClient: redisClient,
	}
}

func (c *URLCodeCounter) Next(ctx context.Context) (uint64, error) {
	callCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	val, err := c.redisClient.Incr(callCtx, redisCounterKey).Result()
	if err != nil {
		return 0, err
	}

	return uint64(val), nil
}
