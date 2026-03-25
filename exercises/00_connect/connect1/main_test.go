package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	client := Connect()
	require.NotNil(t, client)
	defer client.Close()

	result, err := Ping(client)
	require.NoError(t, err)
	assert.Equal(t, "PONG", result)
}
