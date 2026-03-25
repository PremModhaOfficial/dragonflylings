package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewPool(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		PoolSize:     10,
		MinIdleConns: 5,
		PoolTimeout:  2 * time.Second,
	})
}

func Ping(client *redis.Client) error {
	ctx := context.Background()
	return client.Ping(ctx).Err()
}

func PoolStats(client *redis.Client) *redis.PoolStats {
	return client.PoolStats()
}
