package main

// EXERCISE: lists4 - LLEN + LTRIM for Bounded Lists
//
// PREDICT: Before fixing anything, answer:
//   - What does LTRIM do to a list?
//   - What does LTRIM mylist 0 4 do to a 10-element list?
//   - What happens to memory if you push to a list forever without trimming?
//
// The test stores the last 5 activity events. After pushing 10 events,
// the list should still have exactly 5 items (the most recent ones).
// BUG: AddEvent pushes to the list but never trims it — the list grows without bound.
//
// TODO: After pushing, call LTRIM to keep only the last maxItems items.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// AddEvent adds an event to the front of the list and trims to maxItems.
// BUG: Missing the LTRIM call — the list grows unbounded.
func AddEvent(client *redis.Client, listKey, event string, maxItems int64) error {
	ctx := context.Background()
	if err := client.LPush(ctx, listKey, event).Err(); err != nil {
		return err
	}
	// TODO: add client.LTrim(ctx, listKey, 0, maxItems-1).Err()
	client.LTrim(ctx, listKey, 0, maxItems-1)
	// This keeps only the first maxItems elements (indices 0 to maxItems-1)
	return nil
}

// GetRecentEvents returns all events currently in the bounded list.
func GetRecentEvents(client *redis.Client, listKey string) ([]string, error) {
	ctx := context.Background()
	return client.LRange(ctx, listKey, 0, -1).Result()
}

// EventCount returns the current list length.
func EventCount(client *redis.Client, listKey string) (int64, error) {
	ctx := context.Background()
	return client.LLen(ctx, listKey).Result()
}
