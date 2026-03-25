package main

// EXERCISE: prod4 - Session Storage
//
// PREDICT: Before writing any code, answer in your head:
//   Redis hashes are ideal for sessions: one key per session, one field per
//   attribute. But what happens when Redis is down? Should your login flow fail?
//   What should GetSession return when Redis is unreachable?
//
// Production session storage concerns:
//   1. TTL: sessions should expire (use EXPIRE after every HSET)
//   2. TTL refresh: active sessions should stay alive (refresh TTL on access)
//   3. Graceful degradation: if Redis is down, return empty session, not error
//   4. Partial updates: HSET can update individual fields without re-serializing
//
// TODO: Fix THREE bugs:
//   SetSession: must set TTL after storing the hash
//   GetSession: must refresh TTL on access (sliding expiry)
//   GetSession: must return empty session (not error) when Redis is down

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

const sessionTTL = 30 * time.Minute

// Session represents an authenticated user session.
type Session struct {
	ID   string
	Data map[string]string
}

func sessionKey(id string) string {
	return "session:" + id
}

// SetSession stores session data as a Redis hash with a 30-minute TTL.
// BUG: no TTL is set — sessions never expire.
func SetSession(ctx context.Context, client *redis.Client, s Session) error {
	key := sessionKey(s.ID)
	if err := client.HSet(ctx, key, toFields(s.Data)).Err(); err != nil {
		return err
	}
	// BUG: missing client.Expire(ctx, key, sessionTTL) call
	return nil
}

// GetSession retrieves a session by ID and refreshes its TTL.
// BUG 1: does not refresh TTL on access (sessions expire even when active)
// BUG 2: returns error when Redis is unreachable (should return empty session)
func GetSession(ctx context.Context, client *redis.Client, id string) (Session, error) {
	key := sessionKey(id)
	fields, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		// BUG: propagates Redis error instead of gracefully degrading
		return Session{}, err
	}
	// BUG: missing TTL refresh: client.Expire(ctx, key, sessionTTL)
	return Session{ID: id, Data: fields}, nil
}

// UpdateField updates a single session field without touching other fields.
// Returns false if the session doesn't exist.
func UpdateField(ctx context.Context, client *redis.Client, id, field, value string) (bool, error) {
	key := sessionKey(id)
	// Check session exists
	exists, err := client.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return false, err
	}
	if err := client.HSet(ctx, key, field, value).Err(); err != nil {
		return false, err
	}
	// Refresh TTL on update
	client.Expire(ctx, key, sessionTTL)
	return true, nil
}

// DeleteSession removes a session from storage.
func DeleteSession(ctx context.Context, client *redis.Client, id string) error {
	return client.Del(ctx, sessionKey(id)).Err()
}

// toFields converts a map to a flat key-value slice for HSET.
func toFields(data map[string]string) []interface{} {
	fields := make([]interface{}, 0, len(data)*2)
	for k, v := range data {
		fields = append(fields, k, v)
	}
	return fields
}

// isRedisDown returns true if the error indicates Redis is unreachable.
func isRedisDown(err error) bool {
	if err == nil || errors.Is(err, redis.Nil) {
		return false
	}
	return true
}
