package main

// EXERCISE: prod6 - Circuit Breaker
//
// PREDICT: Before writing any code, answer in your head:
//   Dragonfly is degraded — every Redis call times out after 5 seconds.
//   Your service has 100 goroutines, each making Redis calls.
//   What happens over the next 30 seconds without a circuit breaker?
//   (Hint: how many goroutines are waiting? What happens to new requests?)
//
// A circuit breaker has three states:
//   CLOSED  (normal): calls pass through, failures are counted
//   OPEN    (tripped): calls immediately return fallback error (no Redis hit)
//   HALF-OPEN (probing): one call allowed through to test recovery
//
// State transitions:
//   CLOSED → OPEN:      failures >= threshold within window
//   OPEN → HALF-OPEN:   after cooldown period elapses
//   HALF-OPEN → CLOSED: probe call succeeds
//   HALF-OPEN → OPEN:   probe call fails
//
// TODO: Implement the state machine in CircuitBreaker.
// The struct fields and State type are defined. Implement:
//   canCall() — check if circuit allows a call
//   recordSuccess() — handle successful call
//   recordFailure() — handle failed call, potentially open circuit
//   Get() — use canCall/record to wrap Redis calls

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// ErrCircuitOpen is returned when the circuit breaker is open.
var ErrCircuitOpen = errors.New("circuit breaker is open: Redis calls are being short-circuited")

// State represents the circuit breaker state.
type State int

const (
	StateClosed   State = iota // Normal operation
	StateOpen                  // Failing — reject calls immediately
	StateHalfOpen              // Probing — allow one call through
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// CircuitBreaker wraps a Redis client with failure detection and fast-fail behavior.
type CircuitBreaker struct {
	client    *redis.Client
	mu        sync.Mutex
	state     State
	failures  int
	threshold int       // failures before opening
	openUntil time.Time // when to try half-open
	cooldown  time.Duration
}

// NewCircuitBreaker creates a circuit breaker with the given failure threshold and cooldown.
func NewCircuitBreaker(client *redis.Client, threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		client:    client,
		state:     StateClosed,
		threshold: threshold,
		cooldown:  cooldown,
	}
}

// GetState returns the current circuit breaker state (for testing/monitoring).
func (cb *CircuitBreaker) GetState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.currentState()
}

// currentState resolves the effective state, transitioning OPEN → HALF-OPEN when cooldown elapses.
// Must be called with cb.mu held.
func (cb *CircuitBreaker) currentState() State {
	if cb.state == StateOpen && time.Now().After(cb.openUntil) {
		cb.state = StateHalfOpen
	}
	return cb.state
}

// Get retrieves a value from Redis through the circuit breaker.
// BUG: no circuit breaker logic — always calls Redis regardless of failure state.
// When Redis is degraded, every call will time out, goroutines pile up, and
// the service becomes unresponsive.
func (cb *CircuitBreaker) Get(ctx context.Context, key string) (string, error) {
	// BUG: missing canCall() check — should return ErrCircuitOpen when state is Open
	// BUG: missing recordSuccess()/recordFailure() calls
	return cb.client.Get(ctx, key).Result()
}

// canCall returns true if the circuit allows a call through.
// Must be called with cb.mu held.
// TODO: Implement this — return false when state is Open, true otherwise.
func (cb *CircuitBreaker) canCall() bool {
	// BUG: always returns true — circuit never blocks calls
	return true
}

// recordSuccess resets the failure count and closes the circuit.
// Must be called with cb.mu held.
// TODO: Implement this.
func (cb *CircuitBreaker) recordSuccess() {
	// BUG: not implemented
}

// recordFailure increments failures and opens the circuit if threshold is reached.
// Must be called with cb.mu held.
// TODO: Implement this.
func (cb *CircuitBreaker) recordFailure() {
	// BUG: not implemented
}
