package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

func MustConnectRedis() *redis.Client {

	client := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	r := client.Ping(timeout)
	if r.Err() != nil {
		panic(r.Err())
	}
	return client
}
