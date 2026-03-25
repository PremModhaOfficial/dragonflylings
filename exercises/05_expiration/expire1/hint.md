## Hint 1

`SetWithExpiry` calls `client.Expire(ctx, key, 0)`. What happens when you pass `0` as a duration to EXPIRE? Check what the function signature expects vs. what variable holds the requested duration.

## Hint 2

`HasExpiry` uses `ttl == -1` to decide. Run `TTL` on a key and look at what values it returns:
- Positive value: key has TTL, this many seconds left
- `-1`: key exists but has **no** expiry (persistent)
- `-2`: key doesn't exist

Which of these means "has an active expiry"?

## Hint 3

`MakePersistent` calls `client.Del(ctx, key)`. DEL removes the key entirely. You want to *remove the expiry* but keep the value. Look at the Redis command `PERSIST` — in go-redis that's `client.Persist(ctx, key)`.
