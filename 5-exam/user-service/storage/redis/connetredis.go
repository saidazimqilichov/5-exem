package storage

import (
	"context"
	"log"
	"github.com/go-redis/redis/v8"
)

func ConnectRedis() (*redis.Client, error) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis-db:6379",
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Redis")

	return rdb, nil
}
