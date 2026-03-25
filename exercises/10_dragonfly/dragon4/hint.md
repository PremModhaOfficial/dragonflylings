## Hint 1

**WaitForReplication bug:** In Dragonfly standalone mode, `WAIT` always returns `0` because there are no replicas. This is correct behavior — not an error. Remove the check that treats `result == 0` as a failure. The function should return `(0, nil)` in this case.

**GetEncoding bug:** Replace the `panic()` with a graceful return. If the encoding isn't in `knownEncodings`, still return an `EncodingInfo` with `IsKnown: false` — don't crash.

## Hint 2

For `WaitForReplication`, just remove the `if result == 0` block entirely. Return `result, nil` unconditionally (the error from `Wait` is already handled).

For `GetEncoding`, replace the panic:
```go
if !knownEncodings[result] {
    // Return gracefully — Dragonfly may use different encoding names
    return EncodingInfo{Encoding: result, IsKnown: false}, nil
}
```

## Hint 3

Complete fixes:

```go
func WaitForReplication(ctx context.Context, client *redis.Client, numReplicas int, timeoutMs int64) (int64, error) {
    result, err := client.Wait(ctx, numReplicas, time.Duration(timeoutMs)*time.Millisecond).Result()
    if err != nil {
        return 0, err
    }
    return result, nil  // 0 is correct for standalone Dragonfly
}

func GetEncoding(ctx context.Context, client *redis.Client, key string) (EncodingInfo, error) {
    result, err := client.ObjectEncoding(ctx, key).Result()
    if err != nil {
        return EncodingInfo{}, err
    }
    return EncodingInfo{Encoding: result, IsKnown: knownEncodings[result]}, nil
}
```
