package main

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
func SetSession(ctx context.Context, client *redis.Client, s Session) error {
	key := sessionKey(s.ID)
	if err := client.HSet(ctx, key, toFields(s.Data)).Err(); err != nil {
		return err
	}
	return client.Expire(ctx, key, sessionTTL).Err()
}

// GetSession retrieves a session by ID, refreshes its TTL, and returns an
// empty session (not an error) when Redis is unreachable.
func GetSession(ctx context.Context, client *redis.Client, id string) (Session, error) {
	key := sessionKey(id)
	fields, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		if isRedisDown(err) {
			// Graceful degradation: Redis is down, return empty session
			return Session{ID: id, Data: map[string]string{}}, nil
		}
		return Session{}, err
	}
	// Refresh TTL on access (sliding expiry)
	if len(fields) > 0 {
		client.Expire(ctx, key, sessionTTL)
	}
	return Session{ID: id, Data: fields}, nil
}

// UpdateField updates a single session field without touching other fields.
func UpdateField(ctx context.Context, client *redis.Client, id, field, value string) (bool, error) {
	key := sessionKey(id)
	exists, err := client.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return false, err
	}
	if err := client.HSet(ctx, key, field, value).Err(); err != nil {
		return false, err
	}
	client.Expire(ctx, key, sessionTTL)
	return true, nil
}

// DeleteSession removes a session from storage.
func DeleteSession(ctx context.Context, client *redis.Client, id string) error {
	return client.Del(ctx, sessionKey(id)).Err()
}

func toFields(data map[string]string) []interface{} {
	fields := make([]interface{}, 0, len(data)*2)
	for k, v := range data {
		fields = append(fields, k, v)
	}
	return fields
}

func isRedisDown(err error) bool {
	if err == nil || errors.Is(err, redis.Nil) {
		return false
	}
	return true
}
