## Hint 1

`pipe.Exec(ctx)` returns `([]redis.Cmder, error)`. The top-level error is set if ANY command fails, but the slice of `Cmder` contains the result of every individual command -- including which ones succeeded and which failed.

## Hint 2

Each `redis.Cmder` has an `Err() error` method. After `Exec`, loop over the returned commands and call `.Err()` on each. If `.Err() != nil`, that specific command failed.

## Hint 3

```go
cmds, _ := pipe.Exec(ctx)
failures := 0
for _, cmd := range cmds {
    if cmd.Err() != nil {
        failures++
    }
}
return failures, nil
```
