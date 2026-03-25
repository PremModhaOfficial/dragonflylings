package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// atomicTransferScript moves `amount` from one counter to another atomically.
// KEYS[1] = fromKey, KEYS[2] = toKey
// ARGV[1] = amount to transfer
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
	_, err := client.Eval(ctx, atomicTransferScript, []string{fromKey, toKey}, amount).Int()
	return err
}
