package main

// EXERCISE: sets1 - SADD/SMEMBERS/SISMEMBER
//
// PREDICT: Before fixing anything, answer:
//   - What's the key property of a Set that a List doesn't have?
//   - What does SADD return when you add a value that already exists?
//   - Can a Redis Set contain the same value twice?
//
// The test builds a tag system where each tag is unique.
// BUG: AddTag uses SADD but HasTag uses SISMEMBER wrong — it checks membership
//      using the list-oriented LPOS instead of the set-oriented SISMEMBER.
//
// Wait, the real bug: AddTag uses LPUSH (list) instead of SADD (set),
// so duplicates are stored and membership check (SISMEMBER) doesn't work.
//
// TODO: Fix AddTag to use SADD, and GetTags to use SMEMBERS.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// AddTag adds a tag to the item's tag set.
// BUG: Uses LPush (list) — allows duplicate tags and wrong data structure.
func AddTag(client *redis.Client, itemKey, tag string) error {
	ctx := context.Background()
	return client.SAdd(ctx, itemKey+":tags", tag).Err() // TODO: use SAdd
}

// GetTags returns all unique tags for an item.
// BUG: Uses LRange (list) — doesn't deduplicate.
func GetTags(client *redis.Client, itemKey string) ([]string, error) {
	ctx := context.Background()
	return client.SMembers(ctx, itemKey+":tags").Result() // TODO: use SMembers
}

// HasTag checks whether an item has a specific tag.
// BUG: With list storage, SISMEMBER doesn't work (wrong type).
func HasTag(client *redis.Client, itemKey, tag string) (bool, error) {
	ctx := context.Background()
	return client.SIsMember(ctx, itemKey+":tags", tag).Result()
}
