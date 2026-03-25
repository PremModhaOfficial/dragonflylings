package main

import (
	"context"
	"testing"

	"dragonflylings/lib/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetAndGetProfile(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon3")

	original := Profile{
		Name:  "Alice",
		Age:   30,
		Email: "alice@example.com",
		Tags:  []string{"admin", "user"},
	}

	err := SetProfile(ctx, client, key, original)
	require.NoError(t, err, "SetProfile should succeed")

	got, err := GetProfile(ctx, client, key)
	require.NoError(t, err, "GetProfile should succeed")

	assert.Equal(t, original.Name, got.Name)
	assert.Equal(t, original.Age, got.Age)
	assert.Equal(t, original.Email, got.Email)
	assert.Equal(t, original.Tags, got.Tags)
}

func TestGetProfileField_Name(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon3-field")

	p := Profile{Name: "Bob", Age: 25, Email: "bob@example.com", Tags: []string{}}
	require.NoError(t, SetProfile(ctx, client, key, p))

	name, err := GetProfileField(ctx, client, key, "$.name")
	require.NoError(t, err)
	assert.Equal(t, "Bob", name)
}

func TestGetProfileField_Email(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon3-email")

	p := Profile{Name: "Carol", Age: 35, Email: "carol@example.com", Tags: []string{}}
	require.NoError(t, SetProfile(ctx, client, key, p))

	email, err := GetProfileField(ctx, client, key, "$.email")
	require.NoError(t, err)
	assert.Equal(t, "carol@example.com", email)
}

func TestSetProfile_UpdateField(t *testing.T) {
	client := testutil.NewTestClient(t)
	ctx := context.Background()
	key := testutil.UniqueKey("dragon3-update")

	original := Profile{Name: "Dave", Age: 40, Email: "dave@example.com", Tags: []string{"user"}}
	require.NoError(t, SetProfile(ctx, client, key, original))

	// Update just the age using JSON.SET with a path
	err := client.Do(ctx, "JSON.SET", key, "$.age", "41").Err()
	require.NoError(t, err)

	updated, err := GetProfile(ctx, client, key)
	require.NoError(t, err)
	assert.Equal(t, 41, updated.Age, "age should be updated to 41")
	assert.Equal(t, "Dave", updated.Name, "name should be unchanged")
}
