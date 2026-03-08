package cache

import (
	"context"
	"time"

	"github.com/narvdeshwar/IOT-rag/internal/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func Init() {
	Client = redis.NewClient(&redis.Options{
		Addr: config.Load().RedisURL,
	})
}

func Get(ctx context.Context, key string) (string, error) {
	return Client.Get(ctx, key).Result()
}

func Set(ctx context.Context, key, value string, ttl time.Duration) {
	Client.Set(ctx, key, value, ttl)
}
