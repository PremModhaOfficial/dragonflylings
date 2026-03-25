## Hint 1

The state machine is already sketched. You need to implement four things:

1. `canCall()`: return `false` if state is `StateOpen` (after checking cooldown via `currentState()`)
2. `recordSuccess()`: set `cb.state = StateClosed`, `cb.failures = 0`
3. `recordFailure()`: increment `cb.failures`; if `>= cb.threshold`, set `StateOpen` and `cb.openUntil = time.Now().Add(cb.cooldown)`
4. `Get()`: lock, call `canCall()`, unlock, call Redis, lock again, call record, unlock

Note: `redis.Nil` means "key not found" — this is NOT a circuit breaker failure.

## Hint 2

`Get()` structure:

```go
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
```

## Hint 3

Complete state machine implementations:

```go
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
```
