package main

// EXERCISE: lua4 - Dragonfly Hashtag Gotcha
//
// PREDICT: Before writing any code, answer in your head:
//   Dragonfly has multiple shards (like Redis Cluster).
//   If a Lua script touches 2 keys on 2 different shards, what happens?
//   Can the server execute that script atomically across shards?
//
// Dragonfly distributes keys across shards using a hash of the key name.
// A Lua script must only access keys that all hash to the SAME shard.
//
// Solution: hashtags. The {tag} portion of a key name is used for shard
// placement. Keys sharing the same {tag} always land on the same shard.
//
//   "{user:42}:balance"  → shard determined by "user:42"
//   "{user:42}:reserved" → shard determined by "user:42" ← SAME SHARD ✓
//
// TODO: Fix MakeAccountKeys to use hashtag notation so both keys always
// land on the same Dragonfly shard. Wrap the stable part of the key name
// in curly braces: {accountID}

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// reserveScript atomically reserves an amount from a balance.
// If balance >= amount: deducts from balance, adds to reserved. Returns 1.
// If balance < amount: returns 0 (no change).
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
// BUG: keys do not use hashtag notation — in Dragonfly's multi-shard mode
// these keys may land on different shards, causing the Lua script to fail
// with: ERR CROSSSLOT Keys in request don't hash to the same slot
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
