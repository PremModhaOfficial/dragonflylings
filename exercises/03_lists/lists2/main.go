package main

// EXERCISE: lists2 - LRANGE Pagination
//
// PREDICT: Before fixing anything, answer:
//   - What does LRANGE mylist 0 -1 return? // reversed?
//   - What does LRANGE mylist 0 0 return? // nothing?
//   - If a list has 10 items (indices 0-9), what does LRANGE mylist 0 4 return? 3 items ?
//
// The test paginates a list of activity events.
// BUG: GetPage uses 1-based indexing instead of Redis's 0-based indexing.
//      Page 1 should be items at index 0 to pageSize-1, not 1 to pageSize.
//
// TODO: Fix the index calculation to use 0-based page indexing.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// AddActivity prepends an activity to the front of the list (most recent first).
func AddActivity(client *redis.Client, listKey, activity string) error {
	ctx := context.Background()
	return client.LPush(ctx, listKey, activity).Err()
}

// GetPage returns a page of activities. Page 1 = first page.
// BUG: start is calculated as page*pageSize (1-based), should be (page-1)*pageSize (0-based).
func GetPage(client *redis.Client, listKey string, page, pageSize int64) ([]string, error) {
	ctx := context.Background()
	// BUG: 1-based calculation — page 1 starts at index 1, skipping the first item
	start := (page - 1) * pageSize // TODO: should be (page-1) * pageSize
	stop := start + pageSize - 1
	return client.LRange(ctx, listKey, start, stop).Result()
}

// GetAll returns all items in the list.
func GetAll(client *redis.Client, listKey string) ([]string, error) {
	ctx := context.Background()
	return client.LRange(ctx, listKey, 0, -1).Result()
}
