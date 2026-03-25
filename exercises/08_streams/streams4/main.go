package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: streams4 - Crash Recovery with XCLAIM
//
// PREDICT: Before writing any code, answer in your head:
//   What is XPENDING and what does it show?
//   What does XCLAIM do? Who calls it and why?
//   What is "min-idle-time" in XCLAIM?
//
// TODO: Fix the two bugs below.

// FindPendingMessages returns message IDs that are pending (delivered but unacknowledged).
// BUG: Queries XPending summary -- returns aggregate stats, not individual message IDs.
// Use XPendingExt to get per-message details including IDs.
func FindPendingMessages(client *redis.Client, ctx context.Context, stream, group string) ([]string, error) {
	// BUG: XPending returns a summary (*XPending) with Count and bounds,
	// but NOT individual message IDs. Use XPendingExt for detailed listing.
	info, err := client.XPending(ctx, stream, group).Result()
	if err != nil {
		return nil, err
	}
	// Can't get IDs from summary -- returns empty
	_ = info
	return nil, nil
}

// ClaimMessages claims idle messages from crashed consumer and assigns to newConsumer.
// BUG: Uses minIdle=1 hour -- only claims messages idle for over 1 hour.
// In tests (milliseconds), nothing is idle that long. Use 0 to claim all pending.
func ClaimMessages(client *redis.Client, ctx context.Context, stream, group, newConsumer string, msgIDs []string) ([]redis.XMessage, error) {
	if len(msgIDs) == 0 {
		return nil, nil
	}
	return client.XClaim(ctx, &redis.XClaimArgs{
		Stream:   stream,
		Group:    group,
		Consumer: newConsumer,
		MinIdle:  1 * time.Hour, // BUG: too long for tests
		Messages: msgIDs,
	}).Result()
}

func main() {}
