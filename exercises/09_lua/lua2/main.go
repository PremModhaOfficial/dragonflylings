package main

// EXERCISE: lua2 - KEYS vs ARGV
//
// PREDICT: Before writing any code, answer in your head:
//   Redis Cluster and Dragonfly shard keys across multiple nodes/threads.
//   How does the server know WHICH shard to send a Lua script to?
//   Can it look inside the script to figure out what keys it uses?
//
// Redis requires you to declare ALL keys a script will access in the KEYS[]
// array — not hardcode them in the script and not pass them as ARGV[].
// This is because the server needs to know key locations BEFORE executing
// the script, so it can route to the right shard or reject cross-shard access.
//
// TODO: Fix TWO bugs:
//   1. The Lua script uses ARGV[1] and ARGV[2] as key names — move them to KEYS
//   2. The Go call passes keys as extra args (ARGV) — move them to the keys slice

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// atomicTransferScript moves `amount` from one counter to another atomically.
// BUG: keys are accessed as ARGV[1] and ARGV[2], not declared in KEYS[].
// In Dragonfly's multi-shard mode this will fail — the server cannot determine
// shard placement when keys are hidden inside ARGV.
const atomicTransferScript = `
local from_val = tonumber(redis.call('GET', KEYS[1]) or '0')
local amount   = tonumber(ARGV[1])
if from_val < amount then
  return redis.error_reply('insufficient balance')
end
redis.call('DECRBY', KEYS[1], amount)
redis.call('INCRBY', KEYS[2], amount)
return 1
`

// AtomicTransfer moves amount from the "from" key to the "to" key atomically.
// Returns an error if "from" has insufficient balance.
func AtomicTransfer(ctx context.Context, client *redis.Client, fromKey, toKey string, amount int64) error {
	// BUG: fromKey and toKey passed as ARGV (extra args), not as keys
	_, err := client.Eval(ctx, atomicTransferScript, []string{fromKey, toKey}, amount).Int()
	return err
}
