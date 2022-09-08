package cache

import (
	"os"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
)

// RedisConnection func for connect to Redis server.
func RedisConnection() (*redis.Client, error) {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	return redis.NewClient(opt), nil
}
