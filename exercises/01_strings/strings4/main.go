package main

// EXERCISE: strings4 - INCR/DECR Atomicity
//
// PREDICT: Before fixing anything, answer:
//   - What's wrong with doing GET → parse → add 1 → SET in separate commands?
//   - What does "atomic" mean in the context of a database command?
//   - Can two goroutines both INCR the same counter simultaneously? What happens?
//
// The test runs concurrent increments. With GET+SET, values get lost.
// With INCR, every increment is counted correctly.
//
// TODO: Replace the GET+parse+SET pattern with the atomic Incr command.

import (
	"context"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// IncrementCounter atomically increments a counter and returns the new value.
// BUG: This implementation uses GET + parse + SET — three non-atomic operations.
// Under concurrent access, two goroutines can read the same value and both write back,
// causing one increment to be lost (lost update problem).
func IncrementCounter(client *redis.Client, counterKey string) (int64, error) {
	ctx := context.Background()

	// TODO: Replace this entire block with a single client.Incr(ctx, counterKey).Result()
	val, err := client.Get(ctx, counterKey).Result()
	if err == redis.Nil {
		val = "0"
	} else if err != nil {
		return 0, err
	}

	n, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}

	n++
	err = client.Set(ctx, counterKey, n, 0).Err()
	return n, err
}

// GetCounter reads the current value of a counter.
func GetCounter(client *redis.Client, counterKey string) (int64, error) {
	ctx := context.Background()
	return client.Get(ctx, counterKey).Int64()
}
