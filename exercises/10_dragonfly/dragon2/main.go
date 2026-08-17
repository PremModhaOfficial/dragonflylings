package main

// EXERCISE: dragon2 - Forkless Snapshots
//
// PREDICT: Before writing any code, answer in your head:
//   Traditional Redis uses fork() to create a snapshot. Why does fork()
//   cause a 2× memory spike? What does the OS do during fork()?
//   How would Dragonfly avoid this?
//
// Redis snapshot (BGSAVE) works by fork()ing the process. The child writes
// the snapshot while the parent serves requests. Copy-on-write means every
// page modified by the parent is duplicated — at peak, memory doubles.
//
// Dragonfly uses a forkless snapshot: fibers iterate over all key-value
// pairs while the server continues serving requests. Memory usage stays
// nearly flat.
//
// TODO: Fix TWO functions:
//   1. GetUsedMemory: should parse "used_memory" from INFO memory, not use DBSIZE
//   2. WaitForSnapshot: should poll INFO persistence until rdb_bgsave_in_progress=0

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// GetUsedMemory returns the current used_memory value in bytes from INFO memory.
// BUG: uses DBSIZE (key count) instead of used_memory from INFO memory.
// These are completely different metrics!
func GetUsedMemory(ctx context.Context, client *redis.Client) (int64, error) {
	// BUG: DBSIZE returns number of keys, not memory usage in bytes
	result, err := client.Info(ctx, "memory").Result()
	if err != nil {
		return 0, err
	}

	return parseMemoryBytes(result)
}

// TriggerSnapshot starts a background save (BGSAVE) and returns immediately.
func TriggerSnapshot(ctx context.Context, client *redis.Client) error {
	return client.BgSave(ctx).Err()
}

// WaitForSnapshot blocks until the background save completes.
// BUG: returns immediately without checking if BGSAVE actually finished.
// This means callers think the snapshot is done when it may still be running.
func WaitForSnapshot(ctx context.Context, client *redis.Client) error {
	// BUG: no polling — snapshot may still be in progress when this returns
	return pollUntilDone(ctx, client, 100*time.Millisecond)
}

// parseInfoField extracts the value of a named field from an INFO response.
// INFO responses look like: "field:value\r\n"
func parseInfoField(info, field string) (string, bool) {
	for _, line := range strings.Split(info, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, field+":") {
			return strings.TrimSpace(strings.TrimPrefix(line, field+":")), true
		}
	}
	return "", false
}

// parseMemoryBytes parses a "used_memory" value from INFO memory output.
func parseMemoryBytes(info string) (int64, error) {
	val, ok := parseInfoField(info, "used_memory")
	if !ok {
		return 0, nil
	}
	return strconv.ParseInt(val, 10, 64)
}

// pollUntilDone polls INFO persistence until rdb_bgsave_in_progress is 0.
// Use this in WaitForSnapshot.
func pollUntilDone(ctx context.Context, client *redis.Client, interval time.Duration) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
			info, err := client.Info(ctx, "persistence").Result()
			if err != nil {
				return err
			}
			if val, ok := parseInfoField(info, "saving"); ok && val == "0" {
				return nil
			}
			if val, ok := parseInfoField(info, "rdb_bgsave_in_progress"); ok && val == "0" {
				return nil
			}
		}
	}
}
