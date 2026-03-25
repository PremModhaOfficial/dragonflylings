package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrCircuitOpen = errors.New("circuit breaker is open: Redis calls are being short-circuited")

type State int

const (
	StateClosed   State = iota
	StateOpen
	StateHalfOpen
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

type CircuitBreaker struct {
	client    *redis.Client
	mu        sync.Mutex
	state     State
	failures  int
	threshold int
	openUntil time.Time
	cooldown  time.Duration
}

func NewCircuitBreaker(client *redis.Client, threshold int, cooldown time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		client:    client,
		state:     StateClosed,
		threshold: threshold,
		cooldown:  cooldown,
	}
}

func (cb *CircuitBreaker) GetState() State {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.currentState()
}

func (cb *CircuitBreaker) currentState() State {
	if cb.state == StateOpen && time.Now().After(cb.openUntil) {
		cb.state = StateHalfOpen
	}
	return cb.state
}

func (cb *CircuitBreaker) Get(ctx context.Context, key string) (string, error) {
	cb.mu.Lock()
	if !cb.canCall() {
		cb.mu.Unlock()
		return "", ErrCircuitOpen
	}
	cb.mu.Unlock()

	val, err := cb.client.Get(ctx, key).Result()

	cb.mu.Lock()
	defer cb.mu.Unlock()
	if err != nil && err != redis.Nil {
		cb.recordFailure()
		return "", err
	}
	cb.recordSuccess()
	return val, err
}

func (cb *CircuitBreaker) canCall() bool {
	state := cb.currentState()
	return state == StateClosed || state == StateHalfOpen
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.state = StateClosed
	cb.failures = 0
}

func (cb *CircuitBreaker) recordFailure() {
	cb.failures++
	if cb.failures >= cb.threshold {
		cb.state = StateOpen
		cb.openUntil = time.Now().Add(cb.cooldown)
	}
}
