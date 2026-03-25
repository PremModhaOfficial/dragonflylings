package main

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
func SetProfile(ctx context.Context, client *redis.Client, key string, p Profile) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return client.Do(ctx, "JSON.SET", key, "$", string(data)).Err()
}

// GetProfile retrieves a Profile stored as native JSON.
// JSON.GET returns a JSON array wrapping the result: [{"name":...}]
func GetProfile(ctx context.Context, client *redis.Client, key string) (Profile, error) {
	result, err := client.Do(ctx, "JSON.GET", key, "$").Text()
	if err != nil {
		return Profile{}, err
	}
	var profiles []Profile
	if err := json.Unmarshal([]byte(result), &profiles); err != nil {
		return Profile{}, err
	}
	if len(profiles) == 0 {
		return Profile{}, redis.Nil
	}
	return profiles[0], nil
}

// GetProfileField retrieves a single string field from a stored Profile using JSONPath.
func GetProfileField(ctx context.Context, client *redis.Client, key, path string) (string, error) {
	result, err := client.Do(ctx, "JSON.GET", key, path).Text()
	if err != nil {
		return "", err
	}
	var values []string
	if err := json.Unmarshal([]byte(result), &values); err != nil {
		return "", err
	}
	if len(values) == 0 {
		return "", nil
	}
	return values[0], nil
}
