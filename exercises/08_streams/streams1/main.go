package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: streams1 - The Producer and Consumer
//
// PREDICT: Before writing any code, answer in your head:
//   What does "*" mean in XADD? What does "0" mean as an ID?
//   What is the difference between reading from "0" vs "$" in XREAD?
//   Why are stream entry IDs formatted as "timestamp-sequence"?
//
// TODO: Fix the two bugs below.

// AddEvent appends an event to the stream and returns the auto-generated entry ID.
func AddEvent(client *redis.Client, ctx context.Context, stream string, fields map[string]interface{}) (string, error) {
	return client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		// BUG: "0" is a specific ID (maps to "0-0"). The first XAdd works, but
		// subsequent calls fail because IDs must be monotonically increasing.
		// Use "*" to let Dragonfly auto-generate a timestamp-based ID.
		ID:     "0",
		Values: fields,
	}).Result()
}

// ReadAllEvents reads all events from the stream from the very beginning.
// Uses Block: -1 for non-blocking reads (Block: 0 would wait forever for new messages).
func ReadAllEvents(client *redis.Client, ctx context.Context, stream string) ([]redis.XMessage, error) {
	results, err := client.XRead(ctx, &redis.XReadArgs{
		// BUG: "$" means "only messages added AFTER this XREAD call starts".
		// To read existing messages from the beginning, use "0".
		Streams: []string{stream, "$"},
		Count:   100,
		Block:   -1, // non-blocking: return immediately even if no new messages
	}).Result()
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	return results[0].Messages, nil
}

func main() {}
