package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: streams2 - Query the Archive
//
// PREDICT: Before writing any code, answer in your head:
//   What do "-" and "+" mean in XRANGE?
//   What does XRANGE return if start > stop?
//   How does XREVRANGE differ from XRANGE?
//
// TODO: Fix the two bugs below.

// CountEvents returns the total number of entries in the stream.
// BUG: appends ":count" to the stream name -- queries the wrong key.
func CountEvents(client *redis.Client, ctx context.Context, stream string) (int64, error) {
	return client.XLen(ctx, stream).Result()
}

// QueryRange returns all stream entries between startID and stopID (inclusive).
// Use "-" for the oldest entry, "+" for the newest.
// BUG: start and stop are reversed -- XRANGE with "+" before "-" returns nothing.
func QueryRange(client *redis.Client, ctx context.Context, stream, start, stop string) ([]redis.XMessage, error) {
	return client.XRange(ctx, stream, start, stop).Result() // BUG: reversed
}

func main() {}
