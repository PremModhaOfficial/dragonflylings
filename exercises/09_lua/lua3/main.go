package main

// EXERCISE: lua3 - EVALSHA and Script Caching
//
// PREDICT: Before writing any code, answer in your head:
//   Every EVAL call sends the full script text over the network.
//   If your script is 500 bytes and you call it 10,000 times/second,
//   how much extra bandwidth is that? What's the alternative?
//
// EVALSHA lets you send a script once (with SCRIPT LOAD), get back a SHA1
// hash, then call it forever with just the 40-byte hash. The server caches
// the compiled script and runs it directly from the SHA.
//
// TODO: Fix TWO things:
//   1. Add a LoadScript function that loads the script and returns its SHA1
//   2. Fix DecrIfPositive to accept a sha string and use EvalSha instead of Eval

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// decrIfPositiveScript decrements a counter only if it's currently > 0.
// Useful for semaphores, rate limiters, and resource pools.
// KEYS[1] = counter key
// Returns: 1 if decremented, 0 if already at zero or below
const decrIfPositiveScript = `
local current = tonumber(redis.call('GET', KEYS[1]) or '0')
if current > 0 then
  redis.call('DECR', KEYS[1])
  return 1
end
return 0
`

// LoadScript loads the decrIfPositive script into Dragonfly's script cache
// and returns the SHA1 hash for use with EvalSha.
// TODO: Implement this function.
// Use client.ScriptLoad(ctx, decrIfPositiveScript).Result()
func LoadScript(ctx context.Context, client *redis.Client) (string, error) {
	// BUG: not implemented — returns empty SHA, causing EvalSha to fail
	return client.ScriptLoad(ctx, decrIfPositiveScript).Result()
}

// DecrIfPositive decrements the counter at key if it's > 0.
// Returns true if decremented, false if the counter was already at zero.
// TODO: Use client.EvalSha with the provided sha instead of client.Eval.
// Using Eval resends the full script on every call — inefficient at scale.
func DecrIfPositive(ctx context.Context, client *redis.Client, sha, key string) (bool, error) {
	// BUG: ignores sha, sends full script text on every call
	result, err := client.Eval(ctx, decrIfPositiveScript, []string{key}).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
