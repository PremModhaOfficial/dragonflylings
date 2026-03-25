package main

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// GetUsedMemory returns the current used_memory value in bytes from INFO memory.
func GetUsedMemory(ctx context.Context, client *redis.Client) (int64, error) {
	info, err := client.Info(ctx, "memory").Result()
	if err != nil {
		return 0, err
	}
	return parseMemoryBytes(info)
}

// TriggerSnapshot starts a background save (BGSAVE) and returns immediately.
func TriggerSnapshot(ctx context.Context, client *redis.Client) error {
	return client.BgSave(ctx).Err()
}

// WaitForSnapshot blocks until the background save completes.
func WaitForSnapshot(ctx context.Context, client *redis.Client) error {
	return pollUntilDone(ctx, client, 100*time.Millisecond)
}

// parseInfoField extracts the value of a named field from an INFO response.
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
				// Check both fields: Dragonfly uses "saving", Redis compat uses "rdb_bgsave_in_progress"
			saving, hasSaving := parseInfoField(info, "saving")
			rdbInProgress, hasRdb := parseInfoField(info, "rdb_bgsave_in_progress")
			if hasSaving && saving == "0" {
				return nil
			}
			if hasRdb && rdbInProgress == "0" {
				return nil
			}
			// Neither field present — assume not in progress
			if !hasSaving && !hasRdb {
				return nil
			}
		}
	}
}
