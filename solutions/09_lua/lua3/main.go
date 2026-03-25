package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// decrIfPositiveScript decrements a counter only if it's currently > 0.
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
func LoadScript(ctx context.Context, client *redis.Client) (string, error) {
	return client.ScriptLoad(ctx, decrIfPositiveScript).Result()
}

// DecrIfPositive decrements the counter at key if it's > 0.
// Returns true if decremented, false if the counter was already at zero.
// Uses EvalSha with pre-loaded SHA to avoid resending script text on every call.
func DecrIfPositive(ctx context.Context, client *redis.Client, sha, key string) (bool, error) {
	result, err := client.EvalSha(ctx, sha, []string{key}).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
