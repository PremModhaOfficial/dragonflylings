package main

// EXERCISE: hashes4 - Hash vs Multiple Strings
//
// PREDICT: Before fixing anything, answer:
//   - If you store user:1:name, user:1:email, user:1:age as 3 string keys,
//     how do you atomically delete the whole user? How do you list all their fields?
//   - What Redis command gives you ALL fields+values of a hash at once?
//   - Can you partially update a hash without touching other fields?
//
// The test stores a user as multiple string keys (the wrong way) and expects
// it to work as a hash (the right way).
// BUG: StoreUser stores each field as a separate top-level string key.
//      GetUser retrieves them individually too.
//
// TODO: Fix StoreUser to use HSet and GetUser to use HGetAll.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// StoreUser stores user data.
// BUG: Stores each field as a separate string key (scattered across keyspace).
func StoreUser(client *redis.Client, userKey string, data map[string]string) error {
	ctx := context.Background()
	// TODO: use client.HSet(ctx, userKey, data).Err() instead
	for field, value := range data {
		if err := client.Set(ctx, userKey+":"+field, value, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}

// GetUser retrieves all user fields.
// BUG: Cannot use HGetAll because data wasn't stored as a hash.
func GetUser(client *redis.Client, userKey string) (map[string]string, error) {
	ctx := context.Background()
	// TODO: use client.HGetAll(ctx, userKey).Result() instead
	fields := []string{"name", "email", "role", "age"}
	result := make(map[string]string, len(fields))
	for _, f := range fields {
		val, err := client.Get(ctx, userKey+":"+f).Result()
		if err == redis.Nil {
			continue
		}
		if err != nil {
			return nil, err
		}
		result[f] = val
	}
	return result, nil
}

// DeleteUser removes all user data.
// BUG: With string keys, must delete each key manually. With hash, Del(userKey) does it all.
func DeleteUser(client *redis.Client, userKey string) error {
	ctx := context.Background()
	fields := []string{"name", "email", "role", "age"}
	keys := make([]string, len(fields))
	for i, f := range fields {
		keys[i] = userKey + ":" + f
	}
	return client.Del(ctx, keys...).Err()
}
