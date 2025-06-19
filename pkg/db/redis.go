package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

func MustConnectRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{Addr: "cache:6379",
		Username: os.Getenv("REDIS_USER"),
		Password: os.Getenv("REDIS_USER_PASSWORD")})
	r := client.Ping(context.Background())
	if r.Err() != nil {
		panic(r.Err())
	}
	return client
}
