package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: expire2 - Memory Limits
//
// PREDICT: Before writing any code, answer in your head:
//   What happens when Dragonfly runs out of memory? Who decides which key gets evicted?
//   What is the config key to read Dragonfly's memory ceiling?
//   What format does the memory value come back in -- bytes, MB, or a string?
//
// TODO: Fix the two bugs below. The config parameter names are wrong.

// GetMemoryLimit reads the current maxmemory setting from Dragonfly (in bytes).
func GetMemoryLimit(client *redis.Client) (int64, error) {
	ctx := context.Background()
	// BUG: wrong config parameter name -- should be "maxmemory" (no hyphen, no underscore)
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

// SetMemoryLimit sets the maxmemory ceiling in bytes (0 = unlimited).
func SetMemoryLimit(client *redis.Client, bytes int64) error {
	ctx := context.Background()
	// BUG: wrong config parameter name -- should be "maxmemory"
	return client.ConfigSet(ctx, "maxmemory", fmt.Sprintf("%d", bytes)).Err()
}

func main() {}
