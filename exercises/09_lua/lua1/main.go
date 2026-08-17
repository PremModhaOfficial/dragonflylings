package main

// EXERCISE: lua1 - Compare and Swap
//
// PREDICT: Before writing any code, answer in your head:
//   What does "atomic" mean in the context of distributed systems?
//   Why can't we do GET + compare + SET in Go application code safely?
//   What could go wrong between the GET and the SET?
//
// A compare-and-swap (CAS) operation sets a key to a new value ONLY IF
// the current value matches what we expect. This is fundamental to
// optimistic locking and preventing lost updates.
//
// TODO: Fix the Lua script below.
// The ARGV indices are SWAPPED — the script compares against the wrong
// argument and sets the wrong value.
// Convention: ARGV[1] = expected (old value), ARGV[2] = newVal (new value)

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
-- language: lua
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
