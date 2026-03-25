package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// reserveScript atomically reserves an amount from a balance.
// KEYS[1] = balance key, KEYS[2] = reserved key
// ARGV[1] = amount to reserve
const reserveScript = `
local balance  = tonumber(redis.call('GET', KEYS[1]) or '0')
local reserved = tonumber(redis.call('GET', KEYS[2]) or '0')
local amount   = tonumber(ARGV[1])
if balance < amount then
  return 0
end
redis.call('SET', KEYS[1], balance - amount)
redis.call('SET', KEYS[2], reserved + amount)
return 1
`

// MakeAccountKeys returns the balance and reserved keys for an account.
// Uses hashtag notation to ensure both keys land on the same Dragonfly shard,
// which is required for multi-key Lua scripts to work correctly.
func MakeAccountKeys(accountID string) (balanceKey, reservedKey string) {
	balanceKey = fmt.Sprintf("{account:%s}:balance", accountID)
	reservedKey = fmt.Sprintf("{account:%s}:reserved", accountID)
	return
}

// Reserve atomically moves `amount` from the balance key to the reserved key.
// Returns true if reservation succeeded, false if insufficient balance.
func Reserve(ctx context.Context, client *redis.Client, accountID string, amount int64) (bool, error) {
	balanceKey, reservedKey := MakeAccountKeys(accountID)
	result, err := client.Eval(ctx, reserveScript, []string{balanceKey, reservedKey}, amount).Int()
	if err != nil {
		return false, err
	}
	return result == 1, nil
}
