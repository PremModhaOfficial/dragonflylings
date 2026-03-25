package main

// EXERCISE: dragon4 - Dragonfly Gotchas
//
// PREDICT: Before writing any code, answer in your head:
//   Redis's WAIT command blocks until N replicas acknowledge writes.
//   Dragonfly in standalone mode has no replicas. What should WAIT return?
//   Does it block forever? Return an error? Return 0 immediately?
//
// This exercise covers three real-world gotchas when migrating from Redis
// to Dragonfly. Code that works correctly on Redis may behave differently
// on Dragonfly:
//
// Gotcha 1: WAIT returns 0 immediately (no replicas in standalone mode)
// Gotcha 2: OBJECT ENCODING may return different values than Redis
// Gotcha 3: Sentinel is not supported — code that connects to Sentinel will fail
//
// TODO: Fix WaitForReplication — it should NOT block when numReplicas=0 is returned.
// Also fix GetEncoding — it should handle unknown encodings gracefully.

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// WaitForReplication waits for numReplicas replicas to acknowledge all writes,
// with the given timeout in milliseconds.
// Returns the number of replicas that acknowledged (0 in standalone Dragonfly).
//
// BUG: treats 0 replicas as a timeout/failure instead of a normal Dragonfly response.
// In Dragonfly standalone mode, WAIT always returns 0 — this is correct, not an error.
func WaitForReplication(ctx context.Context, client *redis.Client, numReplicas int, timeoutMs int64) (int64, error) {
	result, err := client.Wait(ctx, numReplicas, time.Duration(timeoutMs)*time.Millisecond).Result()
	if err != nil {
		return 0, err
	}
	// BUG: treats 0 as an error — but Dragonfly standalone always returns 0
	if result == 0 {
		return 0, fmt.Errorf("replication timeout: 0 replicas acknowledged (expected %d)", numReplicas)
	}
	return result, nil
}

// EncodingInfo holds information about a key's internal encoding.
type EncodingInfo struct {
	Encoding string
	IsKnown  bool // true if this is a well-known Redis encoding
}

// knownEncodings lists standard Redis encodings.
var knownEncodings = map[string]bool{
	"int": true, "embstr": true, "raw": true,
	"ziplist": true, "listpack": true, "quicklist": true,
	"hashtable": true, "skiplist": true, "intset": true,
	"zipmap": true, "linkedlist": true,
}

// GetEncoding returns the internal encoding of a key using OBJECT ENCODING.
// BUG: panics on unrecognized encoding instead of returning it gracefully.
// Dragonfly may return encodings not in the standard Redis list.
func GetEncoding(ctx context.Context, client *redis.Client, key string) (EncodingInfo, error) {
	result, err := client.ObjectEncoding(ctx, key).Result()
	if err != nil {
		return EncodingInfo{}, err
	}
	if !knownEncodings[result] {
		// BUG: panics on unknown encoding — should return gracefully
		panic(fmt.Sprintf("unknown encoding %q — is this Redis or Dragonfly?", result))
	}
	return EncodingInfo{Encoding: result, IsKnown: true}, nil
}
