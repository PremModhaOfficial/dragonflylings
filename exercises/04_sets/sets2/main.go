package main

// EXERCISE: sets2 - SINTER/SUNION/SDIFF
//
// PREDICT: Before fixing anything, answer:
//   - What does SINTER return for two sets? SUNION? SDIFF?
//   - If alice follows {bob, carol, dave} and bob follows {alice, carol, eve},
//     what does SINTER return?
//   - What does SDIFF(alice_follows, bob_follows) return?
//
// The test finds mutual connections between users.
// BUG: CommonFollows uses SUNION (returns ALL follows from both users — the union).
//      It should use SINTER (returns only follows they have in COMMON).
//
// TODO: Change SUnion → SInter in CommonFollows.

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Follow records that followerID follows targetID.
func Follow(client *redis.Client, followerKey, targetID string) error {
	ctx := context.Background()
	return client.SAdd(ctx, followerKey, targetID).Err()
}

// GetFollows returns everyone a user follows.
func GetFollows(client *redis.Client, followerKey string) ([]string, error) {
	ctx := context.Background()
	return client.SMembers(ctx, followerKey).Result()
}

// CommonFollows returns users that BOTH key1 and key2 follow.
// BUG: Uses SUnion (everyone either follows) instead of SInter (everyone both follow).
func CommonFollows(client *redis.Client, key1, key2 string) ([]string, error) {
	ctx := context.Background()
	return client.SUnion(ctx, key1, key2).Result() // TODO: use SInter
}

// AllFollows returns everyone that either user follows (union).
func AllFollows(client *redis.Client, key1, key2 string) ([]string, error) {
	ctx := context.Background()
	return client.SUnion(ctx, key1, key2).Result()
}

// UniqueFollows returns users that key1 follows but key2 does not.
func UniqueFollows(client *redis.Client, key1, key2 string) ([]string, error) {
	ctx := context.Background()
	return client.SDiff(ctx, key1, key2).Result()
}
