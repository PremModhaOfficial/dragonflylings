## Hint 1

`SCRIPT LOAD` sends your script to the server once. The server compiles it, stores it, and returns a SHA1 hash. From then on, `EVALSHA <sha> numkeys ...` runs the cached version — no script text is sent again.

In go-redis: `client.ScriptLoad(ctx, scriptText)` returns `*redis.StringCmd`. Call `.Result()` to get the SHA string.

## Hint 2

`LoadScript` should look like this skeleton:
```go
func LoadScript(ctx context.Context, client *redis.Client) (string, error) {
    return client.ScriptLoad(ctx, decrIfPositiveScript).Result()
}
```

For `DecrIfPositive`, replace `client.Eval(ctx, decrIfPositiveScript, ...)` with `client.EvalSha(ctx, sha, ...)`. The signature is identical — same keys, same args, just SHA instead of script text.

## Hint 3

Complete fix:

```go
func LoadScript(ctx context.Context, client *redis.Client) (string, error) {
    return client.ScriptLoad(ctx, decrIfPositiveScript).Result()
}

func DecrIfPositive(ctx context.Context, client *redis.Client, sha, key string) (bool, error) {
    result, err := client.EvalSha(ctx, sha, []string{key}).Int()
    if err != nil {
        return false, err
    }
    return result == 1, nil
}
```

In production: if `EvalSha` returns `NOSCRIPT` error, fall back to `Eval` and reload. Scripts can be evicted on server restart.
