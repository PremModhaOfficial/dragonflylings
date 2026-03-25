package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBestPrimitive(t *testing.T) {
	tests := []struct {
		scenario string
		want     string
	}{
		{
			scenario: "batch-100-independent-writes",
			want:     "pipeline",
		},
		{
			scenario: "atomic-debit-credit",
			want:     "transaction",
		},
		{
			scenario: "server-side-conditional-set",
			want:     "lua",
		},
	}

	for _, tt := range tests {
		t.Run(tt.scenario, func(t *testing.T) {
			got := BestPrimitive(tt.scenario)
			assert.Equal(t, tt.want, got,
				"wrong primitive for scenario %q: got %q, want %q", tt.scenario, got, tt.want)
		})
	}
}
