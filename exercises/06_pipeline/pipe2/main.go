package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// EXERCISE: pipe2 - Pipeline Partial Failures
//
// PREDICT: Before writing any code, answer in your head:
//   If one command in a pipeline fails, do the other commands still execute?
//   What does pipe.Exec return when some commands fail?
//   How is this different from a MULTI/EXEC transaction?
//
// TODO: Fix CountPipelineFailures to correctly count individual command errors.

// CountPipelineFailures runs 3 commands: 2 that succeed and 1 that fails.
// It returns how many of the 3 commands resulted in an error.
// Setup: key is pre-set to the string "not-a-number" before this function runs.
func CountPipelineFailures(client *redis.Client, ctx context.Context, key string) (int, error) {
	pipe := client.Pipeline()
	pipe.Set(ctx, key+":a", "good1", 0) // will succeed
	pipe.Incr(ctx, key)                 // will FAIL: "not-a-number" is not an integer
	pipe.Set(ctx, key+":b", "good2", 0) // will succeed

	cmds, _ := pipe.Exec(ctx) // top-level error ignored intentionally

	errCount := 0
	for _, cmd := range cmds {
		if cmd.Err() != nil {
			errCount += 1
		}
	}

	// BUG: not inspecting individual command errors -- always returns 0
	return errCount, nil
}

func main() {}
