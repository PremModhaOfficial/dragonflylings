package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnectReachable(t *testing.T) {
	client, err := Connect("localhost:6380")
	require.NoError(t, err, "should connect successfully to Dragonfly")
	require.NotNil(t, client)
	defer client.Close()
}

func TestConnectUnreachable(t *testing.T) {
	// Port 19999 should have nothing listening
	client, err := Connect("localhost:19999")
	assert.Error(t, err, "should return error when Dragonfly is unreachable")
	assert.Nil(t, client, "should return nil client when connection fails")
}
