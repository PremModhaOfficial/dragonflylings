package main

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// WaitForReplication waits for numReplicas replicas to acknowledge all writes.
// In Dragonfly standalone mode, returns (0, nil) immediately — this is correct.
func WaitForReplication(ctx context.Context, client *redis.Client, numReplicas int, timeoutMs int64) (int64, error) {
	result, err := client.Wait(ctx, numReplicas, time.Duration(timeoutMs)*time.Millisecond).Result()
	if err != nil {
		// Dragonfly standalone may not support WAIT — treat as 0 replicas (correct for standalone)
		if strings.Contains(err.Error(), "unknown command") {
			return 0, nil
		}
		return 0, err
	}
	// 0 replicas is correct for Dragonfly standalone — not an error
	return result, nil
}

// EncodingInfo holds information about a key's internal encoding.
type EncodingInfo struct {
	Encoding string
	IsKnown  bool
}

// knownEncodings lists standard Redis encodings.
var knownEncodings = map[string]bool{
	"int": true, "embstr": true, "raw": true,
	"ziplist": true, "listpack": true, "quicklist": true,
	"hashtable": true, "skiplist": true, "intset": true,
	"zipmap": true, "linkedlist": true,
}

// GetEncoding returns the internal encoding of a key using OBJECT ENCODING.
// Handles Dragonfly-specific encodings gracefully (IsKnown=false, no panic).
// If OBJECT is unsupported (Dragonfly may not implement it), returns a
// graceful "unsupported" result rather than an error.
func GetEncoding(ctx context.Context, client *redis.Client, key string) (EncodingInfo, error) {
	result, err := client.ObjectEncoding(ctx, key).Result()
	if err != nil {
		// Dragonfly may not support OBJECT ENCODING — handle gracefully
		if strings.Contains(err.Error(), "unknown command") {
			return EncodingInfo{Encoding: "unsupported", IsKnown: false}, nil
		}
		return EncodingInfo{}, err
	}
	return EncodingInfo{Encoding: result, IsKnown: knownEncodings[result]}, nil
}
