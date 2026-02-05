package common

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context) *redis.Client {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "mas.internal.io:6379",
		Password: "rds-passwd-123456",
		DB:       7,
	})

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	return rdb
}
