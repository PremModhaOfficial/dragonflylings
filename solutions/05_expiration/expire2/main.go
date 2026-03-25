package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func GetMemoryLimit(client *redis.Client) (int64, error) {
	ctx := context.Background()
	vals, err := client.ConfigGet(ctx, "maxmemory").Result()
	if err != nil {
		return 0, err
	}
	raw, ok := vals["maxmemory"]
	if !ok {
		return 0, fmt.Errorf("maxmemory not found in config response")
	}
	var n int64
	fmt.Sscanf(raw, "%d", &n)
	return n, nil
}

func SetMemoryLimit(client *redis.Client, bytes int64) error {
	ctx := context.Background()
	return client.ConfigSet(ctx, "maxmemory", fmt.Sprintf("%d", bytes)).Err()
}

func main() {}
