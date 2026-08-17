package main

// EXERCISE: dragon3 - Native JSON
//
// PREDICT: Before writing any code, answer in your head:
//   Redis stores everything as strings/bytes. If you want to store JSON,
//   you have two approaches: serialize to string yourself, OR use a JSON
//   data type. What's the advantage of a native JSON type over a plain string?
//
// Dragonfly has built-in JSON support (no module required). You can:
//   JSON.SET key path value    — store JSON at a path ("$" = root)
//   JSON.GET key path          — retrieve JSON at a path
//   JSON.GET key $.field       — retrieve a specific field using JSONPath
//
// The JSON value must be a valid JSON string. Passing a Go struct directly
// will NOT work — you must marshal it first.
//
// TODO: Fix SetProfile and GetProfile:
//   SetProfile: marshal the Profile to JSON before passing to JSON.SET
//   GetProfile: unmarshal the returned JSON array (JSON.GET wraps in [])

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

// Profile represents a user profile stored as native JSON in Dragonfly.
type Profile struct {
	Name  string   `json:"name"`
	Age   int      `json:"age"`
	Email string   `json:"email"`
	Tags  []string `json:"tags"`
}

// SetProfile stores a Profile as native JSON at the given key.
// BUG: passes the Go struct directly — JSON.SET expects a JSON string.
// You must marshal Profile to JSON bytes first.
func SetProfile(ctx context.Context, client *redis.Client, key string, p Profile) error {
	// BUG: profile struct passed directly, not as a JSON string
	jsonProfile, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return client.Do(ctx, "JSON.SET", key, "$", jsonProfile).Err()
}

// GetProfile retrieves a Profile stored as native JSON.
// BUG: JSON.GET returns a JSON array wrapping the result (e.g. [{"name":...}]).
// You must unmarshal the outer array and take the first element.
func GetProfile(ctx context.Context, client *redis.Client, key string) (Profile, error) {
	result, err := client.Do(ctx, "JSON.GET", key, "$").Text()
	if err != nil {
		return Profile{}, err
	}
	// BUG: tries to unmarshal directly into Profile, but result is [{"name":...}]
	var p []Profile
	if err := json.Unmarshal([]byte(result), &p); err != nil {
		return Profile{}, err
	}
	return p[0], nil
}

// GetProfileField retrieves a single string field from a stored Profile using JSONPath.
// path should be a JSONPath expression like "$.name" or "$.email"
func GetProfileField(ctx context.Context, client *redis.Client, key, path string) (string, error) {
	result, err := client.Do(ctx, "JSON.GET", key, path).Text()
	if err != nil {
		return "", err
	}
	// JSON.GET with a path returns an array: ["value"]
	var values []string
	if err := json.Unmarshal([]byte(result), &values); err != nil {
		return "", err
	}
	if len(values) == 0 {
		return "", nil
	}
	return values[0], nil
}
