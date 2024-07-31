package redis

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})
}
