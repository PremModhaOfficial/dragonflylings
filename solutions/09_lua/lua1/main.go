package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// casScript atomically compares and swaps a value.
// KEYS[1] = the key to operate on
// ARGV[1] = expected current value
// ARGV[2] = new value to set if current matches
// Returns: 1 if swapped, 0 if current value didn't match
const casScript = `
local current = redis.call('GET', KEYS[1])
if current == false then current = '' end
if current == ARGV[1] then
  redis.call('SET', KEYS[1], ARGV[2])
  return 1
end
return 0
`

// CompareAndSwap atomically sets key to newVal only if its current value equals expected.
// Returns true if the swap was performed, false if the value didn't match.
func CompareAndSwap(ctx context.Context, client *redis.Client, key, expected, newVal string) (bool, error) {
	result, err := client.Eval(ctx, casScript, []string{key}, expected, newVal).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
