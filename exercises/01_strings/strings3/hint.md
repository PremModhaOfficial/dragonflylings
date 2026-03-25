# Hints for strings3

## Hint 1 — The Concept

`SETNX` (Set if Not eXists) is the primitive behind distributed locking. It guarantees that only one client can successfully set a key — whoever runs SETNX first "wins" and all others are told "no."

The key property: SETNX is atomic. There's no race condition between checking if the key exists and setting it. Redis executes it as a single indivisible operation.

Why does this matter? If you did `if key_exists: return False; else: set_key`, another process could slip in between the check and the set — two processes would both think they acquired the lock.

## Hint 2 — The Specific Issue

`client.Set()` always overwrites. It doesn't check if the key already exists.
`client.SetNX()` sets the key **only if it doesn't already exist** and returns a boolean.

The broken `AcquireLock` uses `client.Set()` and always returns `(true, nil)`. Replace it with `client.SetNX()` and return its boolean result.

```go
// Wrong (always succeeds):
err := client.Set(ctx, lockKey, ownerID, ttl).Err()
return true, err

// Right (returns false if lock held):
return client.SetNX(ctx, lockKey, ownerID, ttl).Result()
```

## Hint 3 — Near Solution

```go
func AcquireLock(client *redis.Client, lockKey, ownerID string, ttl time.Duration) (bool, error) {
    ctx := context.Background()
    return client.SetNX(ctx, lockKey, ownerID, ttl).Result()
}
```
